package goftd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

// NetworkObject An object represents the network (Note: The field level constraints listed here might not cover all the constraints on the field. Additional constraints might exist.)
type NetworkObject struct {
	ReferenceObject
	Description     string `json:"description,omitempty"`
	SubType         string `json:"subType"`
	Value           string `json:"value"`
	IsSystemDefined bool   `json:"isSystemDefined,omitempty"`
	Links           *Links `json:"links,omitempty"`
}

// Reference Returns a reference object
func (n *NetworkObject) Reference() *ReferenceObject {
	r := ReferenceObject{
		ID:      n.ID,
		Name:    n.Name,
		Version: n.Version,
		Type:    n.Type,
	}

	return &r
}

// GetNetworkObjects Get a list of network objects
func (f *FTD) GetNetworkObjects(limit int) ([]*NetworkObject, error) {
	var err error

	filter := make(map[string]string)
	filter["limit"] = strconv.Itoa(limit)

	endpoint := apiNetworksEndpoint
	data, err := f.Get(endpoint, filter)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*NetworkObject `json:"items"`
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	return v.Items, nil
}

// GetNetworkObjectByID Get a network object by ID
func (f *FTD) GetNetworkObjectByID(id string) (*NetworkObject, error) {
	var err error

	endpoint := fmt.Sprintf("%s/%s", apiNetworksEndpoint, id)
	data, err := f.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var v *NetworkObject

	err = json.Unmarshal(data, &v)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	return v, nil
}

func (f *FTD) getNetworkObjectBy(filterString string, limit int) ([]*NetworkObject, error) {
	var err error

	filter := make(map[string]string)
	filter["filter"] = filterString
	filter["limit"] = strconv.Itoa(limit)

	endpoint := apiNetworksEndpoint
	data, err := f.Get(endpoint, filter)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*NetworkObject `json:"items"`
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	return v.Items, nil
}

// CreateNetworkObject Create a new network object
func (f *FTD) CreateNetworkObject(n *NetworkObject, duplicateAction int) error {
	var err error

	n.Type = "networkobject"
	_, err = f.Post(apiNetworksEndpoint, n)
	if err != nil {
		ftdErr := err.(*FTDError)
		//spew.Dump(ftdErr)
		if len(ftdErr.Message) > 0 && (ftdErr.Message[0].Code == "duplicateName" || ftdErr.Message[0].Code == "newInstanceWithDuplicateId") {
			if f.debug {
				glog.Warningf("This is a duplicate\n")
			}
			if duplicateAction == DuplicateActionError {
				return err
			}
		} else {
			if f.debug {
				glog.Errorf("Error: %s\n", err)
			}
			return err
		}
	}

	query := fmt.Sprintf("name:%s", n.Name)
	obj, err := f.getNetworkObjectBy(query, 0)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	var o *NetworkObject
	if len(obj) == 1 {
		o = obj[0]
	} else {
		if f.debug {
			glog.Errorf("Error: length of object is not 1\n")
		}
		return err
	}

	switch duplicateAction {
	case DuplicateActionReplace:
		o.Value = n.Value
		o.SubType = n.SubType

		err = f.UpdateNetworkObject(o)
		if err != nil {
			if f.debug {
				glog.Errorf("Error: %s\n", err)
			}
			return err
		}
	}

	*n = *o
	return nil
}

// CreateNetworkObjectsFromIPs Create Network objects from an array of IP
func (f *FTD) CreateNetworkObjectsFromIPs(ips []string) ([]*NetworkObject, error) {
	var err error
	var retval []*NetworkObject

	os, err := f.GetNetworkObjects(0)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	found := make(map[string]bool)

	for i := range ips {
		for o := range os {
			if ips[i] == os[o].Value && os[o].SubType == "HOST" {
				retval = append(retval, os[o])
				found[ips[i]] = true
				break
			}
		}
	}

	for i := range ips {
		if _, ok := found[ips[i]]; !ok {

			n := new(NetworkObject)
			n.Name = ips[i]
			n.Value = ips[i]
			n.SubType = "HOST"

			err = f.CreateNetworkObject(n, DuplicateActionDoNothing)
			if err != nil {
				if f.debug {
					glog.Errorf("Error: %s\n", err)
				}
				return nil, err
			}
			retval = append(retval, n)
		}
	}

	return retval, nil
}

// DeleteNetworkObject Delete a network object
func (f *FTD) DeleteNetworkObject(n *NetworkObject) error {
	var err error

	err = f.Delete(fmt.Sprintf("%s/%s", apiNetworksEndpoint, n.ID))
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}

// DeleteNetworkObjectByID Delete a network object
func (f *FTD) DeleteNetworkObjectByID(id string) error {
	var err error

	err = f.Delete(fmt.Sprintf("%s/%s", apiNetworksEndpoint, id))
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}

// UpdateNetworkObject Updates a network object
func (f *FTD) UpdateNetworkObject(n *NetworkObject) error {
	var err error

	endpoint := fmt.Sprintf("%s/%s", apiNetworksEndpoint, n.ID)
	data, err := f.Put(endpoint, n)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	err = json.Unmarshal(data, &n)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}
