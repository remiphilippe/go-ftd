package goftd

import (
	"testing"

	"github.com/golang/glog"
)

func TestGetPortObjectGroup(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	obj, err := ftd.GetPortObjectGroups(0)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(obj) < 1 {
		t.Errorf("expecting more than 0 results\n")
	}
}

func TestCreatePortObjectGroup(t *testing.T) {
	var err error

	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	n := new(PortObject)
	n.Name = "testPortObj001"
	n.Port = "56789"

	err = ftd.CreateTCPPortObject(n, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g := new(PortObjectGroup)
	g.Name = "testPortObjGroup001"
	g.Objects = append(g.Objects, n.Reference())

	err = ftd.CreatePortObjectGroup(g, DuplicateActionReplace)
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

	err = ftd.DeletePortObjectGroup(g)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	err = ftd.DeletePortObject(n)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}
}

func TestAddDeletePortToPortGroup(t *testing.T) {
	var err error

	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	p1 := new(PortObject)
	p1.Name = "testPortObj001"
	p1.Port = "56789"

	err = ftd.CreateTCPPortObject(p1, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	p2 := new(PortObject)
	p2.Name = "testPortObj002"
	p2.Port = "45678"

	err = ftd.CreateUDPPortObject(p2, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g := new(PortObjectGroup)
	g.Name = "testPortObjGroup001"
	g.Objects = append(g.Objects, p1.Reference())

	err = ftd.CreatePortObjectGroup(g, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	err = ftd.AddToPortObjectGroup(g, p2)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(g.Objects) < 2 || len(g.Objects) > 2 {
		t.Errorf("object not added\n")
	}

	good := 0
	for i := range g.Objects {
		if g.Objects[i].ID != p1.ID || g.Objects[i].ID != p2.ID {
			break
		}
		good++
	}

	if good == 2 {
		t.Errorf("objects are not the ones we were expecting\n")
	}

	err = ftd.DeleteFromPortObjectGroup(g, p1)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(g.Objects) < 1 || len(g.Objects) > 1 {
		t.Errorf("object not removed\n")
	}

	if g.Objects[0].ID != p2.ID {
		t.Errorf("wrong object removed\n")
	}

	err = ftd.DeletePortObjectGroup(g)
	if err != nil {
		t.Logf("error: %s\n", err)
	}

	err = ftd.DeletePortObject(p1)
	if err != nil {
		t.Logf("error: %s\n", err)
	}

	err = ftd.DeletePortObject(p2)
	if err != nil {
		t.Logf("error: %s\n", err)
	}
}

func TestDuplicatePortObjectGroupDoNothing(t *testing.T) {
	var err error

	ftd, err := initTest()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	p1 := new(PortObject)
	p1.Name = "testPortObj001"
	p1.Port = "56789"

	err = ftd.CreateTCPPortObject(p1, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g := new(PortObjectGroup)
	g.Name = "testPortObjGroup001"
	g.Objects = append(g.Objects, p1.Reference())

	err = ftd.CreatePortObjectGroup(g, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	g1 := new(PortObjectGroup)
	g1.Name = "testPortObjGroup001"
	g1.Objects = append(g1.Objects, p1.Reference())

	err = ftd.CreatePortObjectGroup(g1, DuplicateActionError)
	if err != nil {
		return
	}

	t.Errorf("should have returned an error...\n")

	err = ftd.DeletePortObjectGroup(g)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	err = ftd.DeletePortObject(p1)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	return
}
