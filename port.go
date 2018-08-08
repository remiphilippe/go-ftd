package goftd

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
)

// Port Represents a TCP or UDP port
type Port struct {
	ReferenceObject
	Description     string `json:"description,omitempty"`
	Port            string `json:"port,omitempty"`
	IsSystemDefined bool   `json:"isSystemDefined,omitempty"`
	Links           *Links `json:"links,omitempty"`
	//Paging          *Paging `json:"paging,omitempty"`
}

// Reference Returns a reference object
func (p *Port) Reference() *ReferenceObject {
	r := ReferenceObject{
		ID:      p.ID,
		Name:    p.Name,
		Version: p.Version,
		Type:    p.Type,
	}

	return &r
}

func (f *FTD) getPorts(protocol string) ([]*Port, error) {
	var err error
	var endpoint string

	switch protocol {
	case "TCP":
		endpoint = apiTCPPortsEndpoint
	case "UDP":
		endpoint = apiUDPPortsEndpoint
	}

	data, err := f.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*Port `json:"items"`
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

// GetTCPPorts Get a list of tcp ports
func (f *FTD) GetTCPPorts() ([]*Port, error) {
	return f.getPorts("TCP")
}

// GetUDPPorts Get a list of udp ports
func (f *FTD) GetUDPPorts() ([]*Port, error) {
	return f.getPorts("UDP")
}

func (f *FTD) getPortByID(protocol, id string) (*Port, error) {
	var err error
	var endpoint string

	switch protocol {
	case "TCP":
		endpoint = fmt.Sprintf("%s/%s", apiTCPPortsEndpoint, id)
	case "UDP":
		endpoint = fmt.Sprintf("%s/%s", apiUDPPortsEndpoint, id)
	}

	data, err := f.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var v *Port

	err = json.Unmarshal(data, &v)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}

	return v, nil
}

// GetTCPPortByID Get a tcp port by ID
func (f *FTD) GetTCPPortByID(id string) (*Port, error) {
	return f.getPortByID("TCP", id)
}

// GetUDPPortByID Get a udp port by ID
func (f *FTD) GetUDPPortByID(id string) (*Port, error) {
	return f.getPortByID("UDP", id)
}

func (f *FTD) getPortBy(protocol, filterString string) ([]*Port, error) {
	var err error
	var endpoint string

	switch protocol {
	case "TCP":
		endpoint = apiTCPPortsEndpoint
	case "UDP":
		endpoint = apiUDPPortsEndpoint
	}

	filter := make(map[string]string)
	filter["filter"] = filterString

	data, err := f.Get(endpoint, filter)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*Port `json:"items"`
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

func (f *FTD) createPort(p *Port, duplicateAction int) error {
	var err error
	var protocol string
	var endpoint string

	switch p.Type {
	case TypeTCPPortObject:
		protocol = "TCP"
		endpoint = apiTCPPortsEndpoint
	case TypeUDPPortObject:
		protocol = "UDP"
		endpoint = apiUDPPortsEndpoint
	}

	data, err := f.Post(endpoint, p)
	if err != nil {
		ftdErr := err.(*FTDError)
		//spew.Dump(ftdErr)
		if len(ftdErr.Message) > 0 && (ftdErr.Message[0].Code == "duplicateName" || ftdErr.Message[0].Code == "newInstanceWithDuplicateId") {
			if f.debug {
				glog.Errorf("This is a duplicate\n")
			}
			if duplicateAction == DuplicateActionDoNothing {
				return err
			}
		} else {
			if f.debug {
				glog.Errorf("Error: %s\n", err)
			}
			return err
		}
	}

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
}

// CreateTCPPort Creates a new TCP port
func (f *FTD) CreateTCPPort(p *Port, duplicateAction int) error {
	p.Type = TypeTCPPortObject
	return f.createPort(p, duplicateAction)
}

// CreateUDPPort Creates a new UDP port
func (f *FTD) CreateUDPPort(p *Port, duplicateAction int) error {
	p.Type = TypeUDPPortObject
	return f.createPort(p, duplicateAction)
}

// DeletePort Delete a port
func (f *FTD) DeletePort(p *Port) error {
	var err error
	var endpoint string

	switch p.Type {
	case TypeTCPPortObject:
		endpoint = fmt.Sprintf("%s/%s", apiTCPPortsEndpoint, p.ID)
	case TypeUDPPortObject:
		endpoint = fmt.Sprintf("%s/%s", apiUDPPortsEndpoint, p.ID)
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

// UpdatePort Updates a port
func (f *FTD) UpdatePort(p *Port) error {
	var err error
	var endpoint string

	switch p.Type {
	case TypeTCPPortObject:
		endpoint = fmt.Sprintf("%s/%s", apiTCPPortsEndpoint, p.ID)
	case TypeUDPPortObject:
		endpoint = fmt.Sprintf("%s/%s", apiUDPPortsEndpoint, p.ID)
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
