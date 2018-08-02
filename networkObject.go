package goftd

import (
	"encoding/json"
	"fmt"

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
func (f *FTD) GetNetworkObjects() ([]*NetworkObject, error) {
	var err error

	data, err := f.Get("object/networks", nil)
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

	endpoint := fmt.Sprintf("object/networks/%s", id)
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

func (f *FTD) getNetworkObjectBy(filterString string) ([]*NetworkObject, error) {
	var err error

	filter := make(map[string]string)
	filter["filter"] = filterString

	endpoint := "object/networks"
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
	data, err := f.Post("object/networks", n)
	if err != nil {
		ftdErr := err.(*FTDError)
		//spew.Dump(ftdErr)
		if len(ftdErr.Message) > 0 && (ftdErr.Message[0].Code == "duplicateName" || ftdErr.Message[0].Code == "newInstanceWithDuplicateId") {
			if f.debug {
				glog.Errorf("This is a duplicate\n")
			}
			if duplicateAction == duplicateActionDoNothing {
				return err
			}
		} else {
			if f.debug {
				glog.Errorf("Error: %s\n", err)
			}
			return err
		}
	}

	if duplicateAction == duplicateActionDoNothing {
		err = json.Unmarshal(data, &n)
		if err != nil {
			if f.debug {
				glog.Errorf("Error: %s\n", err)
			}
			return err
		}

		return nil
	} else if duplicateAction == duplicateActionReplace {
		query := fmt.Sprintf("name:%s", n.Name)
		obj, err := f.getNetworkObjectBy(query)
		if err != nil {
			if f.debug {
				glog.Errorf("Error: %s\n", err)
			}
			return err
		}

		if len(obj) == 1 {
			o := obj[0]
			o.Value = n.Value
			o.SubType = n.SubType

			err = f.UpdateNetworkObject(o)
			if err != nil {
				if f.debug {
					glog.Errorf("Error: %s\n", err)
				}
				return err
			}

			*n = *o

			return nil

		}
	}

	return nil
}

// DeleteNetworkObject Delete a network object
func (f *FTD) DeleteNetworkObject(n *NetworkObject) error {
	var err error

	err = f.Delete(fmt.Sprintf("object/networks/%s", n.ID))
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

	endpoint := fmt.Sprintf("object/networks/%s", n.ID)
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
