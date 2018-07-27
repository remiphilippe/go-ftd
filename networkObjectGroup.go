package goftd

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
)

// NetworkObjectGroup Network Object Group
type NetworkObjectGroup struct {
	ReferenceObject
	Description     string             `json:"description,omitempty"`
	IsSystemDefined bool               `json:"isSystemDefined,omitempty"`
	Objects         []*ReferenceObject `json:"objects,omitempty"`
	Links           *Links             `json:"links,omitempty"`
}

// Reference Returns a reference object
func (g *NetworkObjectGroup) Reference() *ReferenceObject {
	r := ReferenceObject{
		ID:      g.ID,
		Name:    g.Name,
		Version: g.Version,
		Type:    g.Type,
	}

	return &r
}

// GetNetworkObjectGroups Get a list of network objects
func (f *FTD) GetNetworkObjectGroups() ([]*NetworkObjectGroup, error) {
	var err error

	data, err := f.Get("object/networkgroups")
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*NetworkObjectGroup `json:"items"`
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

// CreateNetworkObjectGroup Create a new network object
func (f *FTD) CreateNetworkObjectGroup(n *NetworkObjectGroup) error {
	var err error

	n.Type = "networkobjectgroup"
	data, err := f.Post("object/networkgroups", n)
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

// DeleteNetworkObjectGroup Delete a network object
func (f *FTD) DeleteNetworkObjectGroup(n *NetworkObjectGroup) error {
	var err error

	endpoint := fmt.Sprintf("object/networkgroups/%s", n.ID)
	err = f.Delete(endpoint)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}
