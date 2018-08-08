package goftd

import (
	"testing"

	"github.com/golang/glog"
)

func TestGetNetworkObjectGroup(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	obj, err := ftd.GetNetworkObjectGroups(0)
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

	err = ftd.CreateNetworkObject(n, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g := new(NetworkObjectGroup)
	g.Name = "testObjGroup001"
	g.Objects = append(g.Objects, n.Reference())

	err = ftd.CreateNetworkObjectGroup(g, DuplicateActionReplace)
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

	n1 := new(NetworkObject)
	n1.Name = "testObj001"
	n1.SubType = "HOST"
	n1.Value = "1.1.1.1"

	err = ftd.CreateNetworkObject(n1, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	n2 := new(NetworkObject)
	n2.Name = "testObj002"
	n2.SubType = "HOST"
	n2.Value = "2.2.2.2"

	err = ftd.CreateNetworkObject(n2, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g := new(NetworkObjectGroup)
	g.Name = "testObjGroup001"
	g.Objects = append(g.Objects, n1.Reference())

	err = ftd.CreateNetworkObjectGroup(g, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	err = ftd.AddToNetworkObjectGroup(g, n2)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(g.Objects) < 2 || len(g.Objects) > 2 {
		t.Errorf("object not added\n")
	}

	good := 0
	for i := range g.Objects {
		if g.Objects[i].ID != n1.ID || g.Objects[i].ID != n2.ID {
			break
		}
		good++
	}

	if good == 2 {
		t.Errorf("objects are not the ones we were expecting\n")
	}

	err = ftd.DeleteFromNetworkObjectGroup(g, n1)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(g.Objects) < 1 || len(g.Objects) > 1 {
		t.Errorf("object not removed\n")
	}

	if g.Objects[0].ID != n2.ID {
		t.Errorf("wrong object removed\n")
	}

	err = ftd.DeleteNetworkObjectGroup(g)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	err = ftd.DeleteNetworkObject(n1)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	err = ftd.DeleteNetworkObject(n2)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}
}

func TestDuplicateNetworkObjectGroupDoNothing(t *testing.T) {
	var err error

	ftd, err := initTest()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	n := new(NetworkObject)
	n.Name = "testObj001"
	n.SubType = "HOST"
	n.Value = "1.1.1.1"

	err = ftd.CreateNetworkObject(n, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g := new(NetworkObjectGroup)
	g.Name = "testObjGroup001"
	g.Objects = append(g.Objects, n.Reference())

	err = ftd.CreateNetworkObjectGroup(g, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g1 := new(NetworkObjectGroup)
	g1.Name = "testObjGroup001"
	g1.Objects = append(g1.Objects, n.Reference())

	err = ftd.CreateNetworkObjectGroup(g1, DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("should not have returned an error...\n")
		return
	}

	err = ftd.DeleteNetworkObjectGroup(g)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	err = ftd.DeleteNetworkObject(n)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	return
}

func TestCreateNetworkObjectGroupFromIPs(t *testing.T) {
	var err error

	ftd, err := initTest()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	ips1 := []string{
		"1.2.3.4",
		"5.6.7.8",
		"9.10.11.12",
	}

	g, err := ftd.CreateNetworkObjectGroupFromIPs("testAutoCreate", ips1, DuplicateActionError)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if g.ID == "" || g.Version == "" {
		t.Errorf("Invalid ID or Version\n")
	}

	if len(g.Objects) != 3 {
		t.Errorf("Not enough objects... %d\n", len(g.Objects))
	}

	err = ftd.DeleteNetworkObjectGroup(g)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	for i := range g.Objects {
		err = ftd.DeleteNetworkObjectByID(g.Objects[i].ID)
		if err != nil {
			t.Errorf("error: %s\n", err)
		}
	}

}
