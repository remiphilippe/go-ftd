package goftd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

// PortObject Represents a TCP or UDP port
type PortObject struct {
	ReferenceObject
	Description     string `json:"description,omitempty"`
	Port            string `json:"port,omitempty"`
	IsSystemDefined bool   `json:"isSystemDefined,omitempty"`
	Links           *Links `json:"links,omitempty"`
}

// Reference Returns a reference object
func (p *PortObject) Reference() *ReferenceObject {
	r := ReferenceObject{
		ID:      p.ID,
		Name:    p.Name,
		Version: p.Version,
		Type:    p.Type,
	}

	return &r
}

func (f *FTD) getPortObjects(protocol string, limit int) ([]*PortObject, error) {
	var err error
	var endpoint string

	switch protocol {
	case "TCP":
		endpoint = apiTCPPortObjectsEndpoint
	case "UDP":
		endpoint = apiUDPPortObjectsEndpoint
	}

	filter := make(map[string]string)
	filter["limit"] = strconv.Itoa(limit)

	data, err := f.Get(endpoint, filter)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*PortObject `json:"items"`
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

// GetTCPPortObjects Get a list of tcp ports
func (f *FTD) GetTCPPortObjects() ([]*PortObject, error) {
	return f.getPortObjects("TCP", 0)
}

// GetUDPPortObjects Get a list of udp ports
func (f *FTD) GetUDPPortObjects() ([]*PortObject, error) {
	return f.getPortObjects("UDP", 0)
}

func (f *FTD) getPortObjectByID(protocol, id string, limit int) (*PortObject, error) {
	var err error
	var endpoint string

	switch protocol {
	case "TCP":
		endpoint = fmt.Sprintf("%s/%s", apiTCPPortObjectsEndpoint, id)
	case "UDP":
		endpoint = fmt.Sprintf("%s/%s", apiUDPPortObjectsEndpoint, id)
	}

	filter := make(map[string]string)
	filter["limit"] = strconv.Itoa(limit)

	data, err := f.Get(endpoint, filter)
	if err != nil {
		return nil, err
	}

	var v *PortObject

	err = json.Unmarshal(data, &v)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	return v, nil
}

// GetTCPPortObjectByID Get a tcp port by ID
func (f *FTD) GetTCPPortObjectByID(id string) (*PortObject, error) {
	return f.getPortObjectByID("TCP", id, 0)
}

// GetUDPPortObjectByID Get a udp port by ID
func (f *FTD) GetUDPPortObjectByID(id string) (*PortObject, error) {
	return f.getPortObjectByID("UDP", id, 0)
}

func (f *FTD) getPortObjectBy(protocol, filterString string) ([]*PortObject, error) {
	var err error
	var endpoint string

	switch protocol {
	case "TCP":
		endpoint = apiTCPPortObjectsEndpoint
	case "UDP":
		endpoint = apiUDPPortObjectsEndpoint
	}

	filter := make(map[string]string)
	filter["filter"] = filterString

	data, err := f.Get(endpoint, filter)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*PortObject `json:"items"`
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

func (f *FTD) createPortObject(p *PortObject, duplicateAction int) error {
	var err error
	var protocol string
	var endpoint string

	switch p.Type {
	case TypeTCPPortObject:
		protocol = "TCP"
		endpoint = apiTCPPortObjectsEndpoint
	case TypeUDPPortObject:
		protocol = "UDP"
		endpoint = apiUDPPortObjectsEndpoint
	}

	_, err = f.Post(endpoint, p)
	if err != nil {
		ftdErr := err.(*FTDError)

		if len(ftdErr.Message) > 0 && (ftdErr.Message[0].Code == "duplicateName" || ftdErr.Message[0].Code == "newInstanceWithDuplicateId") {
			if f.debug {
				glog.Errorf("This is a duplicate\n")
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

	query := fmt.Sprintf("name:%s", p.Name)
	obj, err := f.getPortObjectBy(protocol, query)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	var o *PortObject
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
		o.Port = p.Port

		err = f.UpdatePortObject(o)
		if err != nil {
			if f.debug {
				glog.Errorf("Error: %s\n", err)
			}
			return err
		}
	}

	*p = *o
	return nil

	//
	//
	/*
		if duplicateAction == DuplicateActionDoNothing {
			err = json.Unmarshal(data, &p)
			if err != nil {
				if f.debug {
					glog.Errorf("Error: %s\n", err)
				}
				return err
			}

			return nil
		} else if duplicateAction == DuplicateActionReplace {
			query := fmt.Sprintf("name:%s", p.Name)

			obj, err := f.getPortBy(protocol, query)
			if err != nil {
				if f.debug {
					glog.Errorf("Error: %s\n", err)
				}
				return err
			}

			if len(obj) == 1 {
				o := obj[0]
				o.Port = p.Port

				err = f.UpdatePort(o)
				if err != nil {
					if f.debug {
						glog.Errorf("Error: %s\n", err)
					}
					return err
				}

				*p = *o

				return nil

			}
		}

		return nil
	*/
}

// CreateTCPPortObject Creates a new TCP port
func (f *FTD) CreateTCPPortObject(p *PortObject, duplicateAction int) error {
	p.Type = TypeTCPPortObject
	return f.createPortObject(p, duplicateAction)
}

// CreateUDPPortObject Creates a new UDP port
func (f *FTD) CreateUDPPortObject(p *PortObject, duplicateAction int) error {
	p.Type = TypeUDPPortObject
	return f.createPortObject(p, duplicateAction)
}

// DeletePortObject Delete a port
func (f *FTD) DeletePortObject(p *PortObject) error {
	var err error
	var endpoint string

	switch p.Type {
	case TypeTCPPortObject:
		endpoint = fmt.Sprintf("%s/%s", apiTCPPortObjectsEndpoint, p.ID)
	case TypeUDPPortObject:
		endpoint = fmt.Sprintf("%s/%s", apiUDPPortObjectsEndpoint, p.ID)
	}

	err = f.Delete(endpoint)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}

// UpdatePortObject Updates a port
func (f *FTD) UpdatePortObject(p *PortObject) error {
	var err error
	var endpoint string

	switch p.Type {
	case TypeTCPPortObject:
		endpoint = fmt.Sprintf("%s/%s", apiTCPPortObjectsEndpoint, p.ID)
	case TypeUDPPortObject:
		endpoint = fmt.Sprintf("%s/%s", apiUDPPortObjectsEndpoint, p.ID)
	}

	data, err := f.Put(endpoint, p)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	err = json.Unmarshal(data, &p)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}
