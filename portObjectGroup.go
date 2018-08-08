package goftd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

// PortObjectGroup Port Object Group
type PortObjectGroup struct {
	ReferenceObject
	Description     string             `json:"description,omitempty"`
	IsSystemDefined bool               `json:"isSystemDefined,omitempty"`
	Objects         []*ReferenceObject `json:"objects,omitempty"`
	Links           *Links             `json:"links,omitempty"`
}

// Reference Returns a reference object
func (p *PortObjectGroup) Reference() *ReferenceObject {
	r := ReferenceObject{
		ID:      p.ID,
		Name:    p.Name,
		Version: p.Version,
		Type:    p.Type,
	}

	return &r
}

// GetPortObjectGroups Get all the port object groups within the limit specified
func (f *FTD) GetPortObjectGroups(limit int) ([]*PortObjectGroup, error) {
	var err error

	filter := make(map[string]string)
	filter["limit"] = strconv.Itoa(limit)

	endpoint := apiPortObjectGroupsEndpoint
	data, err := f.Get(endpoint, filter)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*PortObjectGroup `json:"items"`
		//Paging *Paging       `json:"paging"`
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

func (f *FTD) getPortObjectGroupBy(filterString string) ([]*PortObjectGroup, error) {
	var err error

	filter := make(map[string]string)
	filter["filter"] = filterString

	endpoint := apiPortObjectGroupsEndpoint
	data, err := f.Get(endpoint, filter)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*PortObjectGroup `json:"items"`
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

// CreatePortObjectGroup Create a new port object group
func (f *FTD) CreatePortObjectGroup(g *PortObjectGroup, duplicateAction int) error {
	var err error

	g.Type = "portobjectgroup"
	endpoint := apiPortObjectGroupsEndpoint
	_, err = f.Post(endpoint, g)
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

	query := fmt.Sprintf("name:%s", g.Name)
	obj, err := f.getPortObjectGroupBy(query)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	var o *PortObjectGroup
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
		o.Objects = g.Objects

		err = f.UpdatePortObjectGroup(o)
		if err != nil {
			if f.debug {
				glog.Errorf("Error: %s\n", err)
			}
			return err
		}
	}

	*g = *o
	return nil
}

// DeletePortObjectGroup Delete a port object group
func (f *FTD) DeletePortObjectGroup(g *PortObjectGroup) error {
	var err error

	endpoint := fmt.Sprintf("%s/%s", apiPortObjectGroupsEndpoint, g.ID)
	err = f.Delete(endpoint)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}

// UpdatePortObjectGroup Updates a port object group
func (f *FTD) UpdatePortObjectGroup(g *PortObjectGroup) error {
	var err error

	endpoint := fmt.Sprintf("%s/%s", apiPortObjectGroupsEndpoint, g.ID)
	data, err := f.Put(endpoint, g)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	err = json.Unmarshal(data, &g)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}

// AddToPortObjectGroup Add a Port to an Object Group
func (f *FTD) AddToPortObjectGroup(g *PortObjectGroup, p *PortObject) error {
	var err error
	for k := range g.Objects {
		if g.Objects[k].ID == p.ID {
			if f.debug {
				glog.Errorf("object already in object group\n")
				return fmt.Errorf("object already in object group")
			}
		}
	}

	g.Objects = append(g.Objects, p.Reference())

	err = f.UpdatePortObjectGroup(g)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}

// DeleteFromPortObjectGroup Deletes a Port from an Object Group
func (f *FTD) DeleteFromPortObjectGroup(g *PortObjectGroup, p *PortObject) error {
	var err error
	for k := range g.Objects {
		if g.Objects[k].ID == p.ID {
			g.Objects = append(g.Objects[:k], g.Objects[k+1:]...)
			break
		}
	}

	err = f.UpdatePortObjectGroup(g)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}
