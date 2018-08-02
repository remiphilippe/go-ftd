package goftd

import (
	"testing"

	"github.com/golang/glog"
)

func TestGetNetworkObject(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	_, err = ftd.GetNetworkObjects()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	obj2, err := ftd.getNetworkObjectBy("name:any-ipv4")
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(obj2) < 1 || len(obj2) > 1 {
		t.Errorf("expecting %d results, got %d\n", 1, len(obj2))
	}

	if obj2[0].Name != "any-ipv4" {
		t.Errorf("expecting any-ipv4 got %s\n", obj2[0].Name)
	}
}

func TestCreateNetworkObject(t *testing.T) {
	var err error

	ftd, err := initTest()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	/*
		{
		  "version": "string",
		  "name": "string",
		  "description": "string",
		  "subType": "HOST",
		  "value": "string",
		  "isSystemDefined": true,
		  "id": "string",
		  "type": "networkobject"
		}
	*/

	n := new(NetworkObject)
	n.Name = "testObj001"
	n.SubType = "HOST"
	n.Value = "1.1.1.1"

	err = ftd.CreateNetworkObject(n, DuplicateActionReplace)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	if n.ID == "" || n.Version == "" {
		t.Errorf("ID of value is not populated correctly\n")
	}

	err = ftd.DeleteNetworkObject(n)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}
}

func TestDuplicateNetworkObjectDoNothing(t *testing.T) {
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
		glog.Errorf("error: %s\n", err)
		return
	}

	n1 := new(NetworkObject)
	n1.Name = "testObj001"
	n1.SubType = "HOST"
	n1.Value = "1.1.1.1"

	err = ftd.CreateNetworkObject(n1, DuplicateActionDoNothing)
	if err != nil {
		//glog.Errorf("error: %s\n", err)
		return
	}

	t.Errorf("should have returned an error...\n")

	err = ftd.DeleteNetworkObject(n)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	return
}

func TestDuplicateNetworkID(t *testing.T) {
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
		glog.Errorf("error: %s\n", err)
		return
	}

	err = ftd.CreateNetworkObject(n, DuplicateActionDoNothing)
	if err != nil {
		//glog.Errorf("error: %s\n", err)
		return
	}

	t.Errorf("should have returned an error...\n")

	err = ftd.DeleteNetworkObject(n)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	return
}

func TestDuplicateNetworkObjectReplace(t *testing.T) {
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
		glog.Errorf("error: %s\n", err)
		return
	}

	newValue := "1.1.1.2"
	n1 := new(NetworkObject)
	n1.Name = "testObj001"
	n1.SubType = "HOST"
	n1.Value = newValue

	err = ftd.CreateNetworkObject(n1, DuplicateActionReplace)
	if err != nil {
		//glog.Errorf("error: %s\n", err)
		return
	}

	if n1.ID != n.ID {
		t.Errorf("Error ID is different, expecting: %s, got: %s\n", n.ID, n1.ID)
	}

	if n1.Value != newValue {
		t.Errorf("Error Value is not changed, expecting: %s, got: %s\n", newValue, n1.Value)
	}

	err = ftd.DeleteNetworkObject(n)
	if err != nil {
		glog.Errorf("error: %s\n", err)
	}

	return
}
