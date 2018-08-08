package goftd

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
)

/*
{
  "items": [
    {
      "version": "string",
      "name": "string",
      "accessRuleIds": [
        0
      ],
      "defaultAction": {
        "action": "PERMIT",
        "eventLogAction": "LOG_FLOW_START",
        "intrusionPolicy": {
          "id": "string",
          "type": "string",
          "version": "string",
          "name": "string"
        },
        "syslogServer": {
          "id": "string",
          "type": "string",
          "version": "string",
          "name": "string"
        },
        "type": "accessdefaultaction"
      },
      "sslPolicy": {
        "id": "string",
        "type": "string",
        "version": "string",
        "name": "string"
      },
      "id": "string",
      "rules": [
        {
          "id": "string",
          "type": "string",
          "version": "string",
          "name": "string"
        }
      ],
      "identityPolicySetting": {
        "id": "string",
        "type": "string",
        "version": "string",
        "name": "string"
      },
      "securityIntelligence": {
        "id": "string",
        "type": "string",
        "version": "string",
        "name": "string"
      },
      "type": "accesspolicy",
      "links": {
        "self": "string"
      }
    }
  ],
  "paging": {
    "prev": [
      "string"
    ],
    "next": [
      "string"
    ],
    "limit": 0,
    "offset": 0,
    "count": 0,
    "pages": 0
  }
}
*/

// AccessPolicy Access Policy Object
type AccessPolicy struct {
	ReferenceObject
	AccessRuleIDs []int `json:"accessRuleIDs,omitempty"`
	DefaultAction struct {
		Action          string
		EventLogAction  string
		IntrusionPolicy *ReferenceObject `json:"intrusionPolicy,omitempty"`
		SyslogServer    *ReferenceObject `json:"syslogServer,omitempty"`
		Type            string
	}
	SSLPolicy             *ReferenceObject   `json:"sslPolicy,omitempty"`
	Rules                 []*ReferenceObject `json:"rules,omitempty"`
	IdentityPolicySetting *ReferenceObject   `json:"identityPolicySetting,omitempty"`
	SecurityIntelligence  *ReferenceObject   `json:"securityIntelligence,omitempty"`
	Links                 *Links             `json:"links,omitempty"`
	Paging                *Paging            `json:"paging,omitempty"`
}

// Reference Returns a reference object
func (a *AccessPolicy) Reference() *ReferenceObject {
	r := ReferenceObject{
		ID:      a.ID,
		Name:    a.Name,
		Version: a.Version,
		Type:    a.Type,
	}

	return &r
}

// GetAccessPolicies Get a list of access policies
func (f *FTD) GetAccessPolicies() ([]*AccessPolicy, error) {
	var err error

	data, err := f.Get("policy/accesspolicies", nil)
	if err != nil {
		return nil, err
	}

	var v struct {
		Items []*AccessPolicy `json:"items"`
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

// ModifyAccessPolicy Modify access policy
func (f *FTD) ModifyAccessPolicy(n *AccessPolicy, policy string) error {
	var err error

	// Define expected type for this object
	n.Type = "accesspolicy"
	n.DefaultAction.Type = "accessdefaultaction"

	endpoint := fmt.Sprintf("policy/accesspolicies/%s", policy)
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
