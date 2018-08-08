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

	ports, err := ftd.GetTCPPorts()
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

	p := new(Port)
	p.Name = "testPort123"
	p.Port = "123"

	err = ftd.CreateTCPPort(p, DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if p.ID == "" || p.Version == "" {
		t.Errorf("ID of value is not populated correctly\n")
		return
	}

	t.Logf("object p: %+v\n", p)

	p2, err := ftd.getPortBy("TCP", fmt.Sprintf("name:%s", p.Name))
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(p2) == 1 {
		if p2[0].ID != p.ID {
			t.Errorf("expected ID %s, got %s\n", p.ID, p2[0].ID)
		}
	} else {
		t.Errorf("unexpected count for p2: %d\n", len(p2))
		return
	}

	t.Logf("object p2: %+v\n", p2[0])

	err = ftd.DeletePort(p)
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

	p := new(Port)
	p.Name = "testPort123"
	p.Port = "1234"

	err = ftd.CreateUDPPort(p, DuplicateActionDoNothing)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if p.ID == "" || p.Version == "" {
		t.Errorf("ID of value is not populated correctly\n")
		return
	}

	t.Logf("object p: %+v\n", p)

	p2, err := ftd.getPortBy("UDP", fmt.Sprintf("name:%s", p.Name))
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if len(p2) == 1 {
		if p2[0].ID != p.ID {
			t.Errorf("expected ID %s, got %s\n", p.ID, p2[0].ID)
		}
	} else {
		t.Errorf("unexpected count for p2: %d\n", len(p2))
		return
	}

	t.Logf("object p2: %+v\n", p2[0])

	err = ftd.DeletePort(p)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

}
