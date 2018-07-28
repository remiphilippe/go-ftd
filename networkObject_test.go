package goftd

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
)

func TestGetNetworkObject(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	_, err = ftd.GetNetworkObjects()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	//spew.Dump(obj)

	fmt.Printf("Running Query\n")
	obj2, err := ftd.getNetworkObjectBy("name:any-ipv4")
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}
	spew.Dump(obj2)
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

	spew.Dump(n)

	err = ftd.CreateNetworkObject(n, true)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	fmt.Printf("Creating...\n")
	spew.Dump(n)

	err = ftd.DeleteNetworkObject(n)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}
}

func TestDuplicateNetworkObject(t *testing.T) {
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

	spew.Dump(n)

	err = ftd.CreateNetworkObject(n, true)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	n1 := new(NetworkObject)
	n1.Name = "testObj001"
	n1.SubType = "HOST"
	n1.Value = "1.1.1.1"

	spew.Dump(n)

	err = ftd.CreateNetworkObject(n1, true)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	spew.Dump(n)

	err = ftd.DeleteNetworkObject(n)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}
}
