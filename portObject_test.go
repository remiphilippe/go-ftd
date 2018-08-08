package goftd

import (
	"fmt"
	"testing"
)

func TestGetTCPPort(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	ports, err := ftd.GetTCPPortObjects()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(ports) == 0 {
		t.Errorf("ports length is 0\n")
	}
}

func TestTCPPort(t *testing.T) {
	var err error

	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	p1 := new(PortObject)
	p1.Name = "testPort123"
	p1.Port = "123"

	err = ftd.CreateTCPPortObject(p1, DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if p1.ID == "" || p1.Version == "" {
		t.Errorf("ID of value is not populated correctly\n")
		return
	}

	t.Logf("object p1: %+v\n", p1)

	p2, err := ftd.getPortObjectBy("TCP", fmt.Sprintf("name:%s", p1.Name))
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(p2) == 1 {
		if p2[0].ID != p1.ID {
			t.Errorf("expected ID %s, got %s\n", p1.ID, p2[0].ID)
		}
	} else {
		t.Errorf("unexpected count for p2: %d\n", len(p2))
		return
	}

	t.Logf("object p2: %+v\n", p2[0])

	err = ftd.DeletePortObject(p1)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

}

func TestUDPPort(t *testing.T) {
	var err error

	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	p1 := new(PortObject)
	p1.Name = "testPort123"
	p1.Port = "1234"

	err = ftd.CreateUDPPortObject(p1, DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if p1.ID == "" || p1.Version == "" {
		t.Errorf("ID of value is not populated correctly\n")
		return
	}

	t.Logf("object p1: %+v\n", p1)

	p2, err := ftd.getPortObjectBy("UDP", fmt.Sprintf("name:%s", p1.Name))
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(p2) == 1 {
		if p2[0].ID != p1.ID {
			t.Errorf("expected ID %s, got %s\n", p1.ID, p2[0].ID)
		}
	} else {
		t.Errorf("unexpected count for p2: %d\n", len(p2))
		return
	}

	t.Logf("object p2: %+v\n", p2[0])

	err = ftd.DeletePortObject(p1)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

}
