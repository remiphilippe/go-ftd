package goftd

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
)

func TestGetNetworkObjectGroup(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	obj, err := ftd.GetNetworkObjectGroups()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	spew.Dump(obj)
}

func TestCreateNetworkObjectGroup(t *testing.T) {
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

	err = ftd.CreateNetworkObject(n)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	g := new(NetworkObjectGroup)
	g.Name = "testObjGroup001"
	g.Objects = append(g.Objects, n.Reference())

	spew.Dump(g)

	err = ftd.CreateNetworkObjectGroup(g)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	spew.Dump(g)

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
}
