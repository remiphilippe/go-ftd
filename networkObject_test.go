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

	_, err = ftd.GetNetworkObjects(0)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	obj2, err := ftd.getNetworkObjectBy("name:any-ipv4", 0)
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
		t.Errorf("should not have error'd... error: %s\n", err)
		return
	}

	// fmt.Printf("n1\n")
	// spew.Dump(n1)

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

	err = ftd.CreateNetworkObject(n, DuplicateActionError)
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

func TestCreateNetworkObjectFromIPs(t *testing.T) {
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

	ips2 := []string{
		"1.2.3.4",
		"5.6.7.8",
		"9.10.11.12",
		"13.14.15.16",
	}

	ns, err := ftd.CreateNetworkObjectsFromIPs(ips1)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(ns) != 3 {
		t.Errorf("we should have at least 3 members, have %d\n", len(ns))
	}

	ns2, err := ftd.CreateNetworkObjectsFromIPs(ips2)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(ns2) != 4 {
		t.Errorf("we should have at least 4 members, have %d\n", len(ns2))
	}

	found := false
	for i := range ns2 {
		if ns2[i].Value == "13.14.15.16" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("didn't find 13.14.15.16 in array\n")
	}

	for i := range ns2 {
		err = ftd.DeleteNetworkObject(ns2[i])
		if err != nil {
			t.Errorf("error: %s\n", err)
		}
	}

}
