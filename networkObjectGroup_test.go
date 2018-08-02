package goftd

import (
	"testing"
)

func TestGetNetworkObjectGroup(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	obj, err := ftd.GetNetworkObjectGroups()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(obj) < 1 {
		t.Errorf("expecting more than 0 results\n")
	}
}

func TestCreateNetworkObjectGroup(t *testing.T) {
	var err error

	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	n := new(NetworkObject)
	n.Name = "testObj001"
	n.SubType = "HOST"
	n.Value = "1.1.1.1"

	err = ftd.CreateNetworkObject(n, duplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g := new(NetworkObjectGroup)
	g.Name = "testObjGroup001"
	g.Objects = append(g.Objects, n.Reference())

	err = ftd.CreateNetworkObjectGroup(g)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if g.ID == "" || g.Version == "" {
		t.Errorf("ID of value is not populated correctly\n")
	}

	if len(g.Objects) < 1 || len(g.Objects) > 1 {
		t.Errorf("objects not populated\n")
	}

	if g.Objects[0].ID != n.ID {
		t.Errorf("object is not the one we were expecting\n")
	}

	err = ftd.DeleteNetworkObjectGroup(g)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	err = ftd.DeleteNetworkObject(n)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}
}

func TestAddDeleteNetworkToNetworkGroup(t *testing.T) {
	var err error

	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	n := new(NetworkObject)
	n.Name = "testObj001"
	n.SubType = "HOST"
	n.Value = "1.1.1.1"

	err = ftd.CreateNetworkObject(n, duplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	n1 := new(NetworkObject)
	n1.Name = "testObj002"
	n1.SubType = "HOST"
	n1.Value = "2.2.2.2"

	err = ftd.CreateNetworkObject(n1, duplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g := new(NetworkObjectGroup)
	g.Name = "testObjGroup001"
	g.Objects = append(g.Objects, n.Reference())

	err = ftd.CreateNetworkObjectGroup(g)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	err = ftd.AddToNetworkObjectGroup(g, n1)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(g.Objects) < 2 || len(g.Objects) > 2 {
		t.Errorf("object not added\n")
	}

	good := 0
	for i := range g.Objects {
		if g.Objects[i].ID != n.ID || g.Objects[i].ID != n1.ID {
			break
		}
		good++
	}

	if good == 2 {
		t.Errorf("objects are not the ones we were expecting\n")
	}

	err = ftd.DeleteFromNetworkObjectGroup(g, n)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(g.Objects) < 1 || len(g.Objects) > 1 {
		t.Errorf("object not removed\n")
	}

	if g.Objects[0].ID != n1.ID {
		t.Errorf("wrong object removed\n")
	}

	err = ftd.DeleteNetworkObjectGroup(g)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	err = ftd.DeleteNetworkObject(n)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}
}
