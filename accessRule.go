package goftd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

// AccessRule Access Rule Object
type AccessRule struct {
	ReferenceObject
	RuleID              int                `json:"ruleId,omitempty"`
	SourceZones         []*ReferenceObject `json:"sourceZones,omitempty"`
	DestinationZones    []*ReferenceObject `json:"destinationZones,omitempty"`
	SourceNetworks      []*ReferenceObject `json:"sourceNetworks,omitempty"`
	DestinationNetworks []*ReferenceObject `json:"destinationNetworks,omitempty"`
	SourcePorts         []*ReferenceObject `json:"sourcePorts,omitempty"`
	DestinationPorts    []*ReferenceObject `json:"destinationPorts,omitempty"`
	RuleAction          string             `json:"ruleAction,omitempty"`
	EventLogAction      string             `json:"eventLogAction,omitempty"`
	VLANTags            []*ReferenceObject `json:"vlanTags,omitempty"`
	Users               []*ReferenceObject `json:"users,omitempty"`
	IntrusionPolicy     *ReferenceObject   `json:"intrusionPolicy,omitempty"`
	FilePolicy          *ReferenceObject   `json:"filePolicy,omitempty"`
	LogFiles            bool               `json:"logFiles,omitempty"`
	SyslogServer        *ReferenceObject   `json:"syslogServer,omitempty"`
	Links               *Links             `json:"links,omitempty"`
	parent              string
}

// Reference Returns a reference object
func (a *AccessRule) Reference() *ReferenceObject {
	r := ReferenceObject{
		ID:      a.ID,
		Name:    a.Name,
		Version: a.Version,
		Type:    a.Type,
	}

	return &r
}

// GetAccessRules Get a list of access rules
func (f *FTD) GetAccessRules(policy string, limit int) ([]*AccessRule, error) {
	var err error

	filter := make(map[string]string)
	filter["limit"] = strconv.Itoa(limit)

	endpoint := fmt.Sprintf("policy/accesspolicies/%s/accessrules", policy)
	data, err := f.Get(endpoint, filter)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*AccessRule `json:"items"`
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

func (f *FTD) getAccessRuleBy(filterString, policy string) ([]*AccessRule, error) {
	var err error

	filter := make(map[string]string)
	filter["filter"] = filterString

	endpoint := fmt.Sprintf("policy/accesspolicies/%s/accessrules", policy)
	data, err := f.Get(endpoint, filter)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*AccessRule `json:"items"`
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

// CreateAccessRule Create a new access rule
func (f *FTD) CreateAccessRule(n *AccessRule, policy string) error {
	var err error

	// Define expected type for this object
	n.Type = "accessrule"

	endpoint := fmt.Sprintf("policy/accesspolicies/%s/accessrules", policy)
	data, err := f.Post(endpoint, n)
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

	n.parent = policy

	return nil
}

// DeleteAccessRule Delete an access rule
func (f *FTD) DeleteAccessRule(n *AccessRule) error {
	var err error

	endpoint := fmt.Sprintf("policy/accesspolicies/%s/accessrules/%s", n.parent, n.ID)
	err = f.Delete(endpoint)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	return nil
}
