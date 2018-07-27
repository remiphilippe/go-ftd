package goftd

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
)

func initTest() (*FTD, error) {
	params := make(map[string]string)
	params["grant_type"] = "password"
	params["username"] = os.Getenv("FTD_USER")
	params["password"] = os.Getenv("FTD_PASSWORD")
	params["debug"] = "true"

	ftd, err := NewFTD(os.Getenv("FTD_HOST"), params)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return nil, err
	}

	return ftd, nil
}

func TestToken(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	spew.Dump(ftd)
}

func TestGetNetworkObject(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	obj, err := ftd.GetNetworkObjects()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	spew.Dump(obj)
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

	err = ftd.CreateNetworkObject(n)
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
