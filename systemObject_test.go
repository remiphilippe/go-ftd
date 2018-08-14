package goftd

import (
	"testing"
)

func TestGetNetworkAny(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	obj, err := ftd.GetNetworkAny()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if obj.ID == "" {
		t.Errorf("error: empty object ID\n")
		return
	}

	if obj.Name != "0.0.0.0" {
		t.Errorf("error: wrong object name %s\n", obj.Name)
		return
	}
}
