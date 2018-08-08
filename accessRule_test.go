package goftd

import (
	"testing"
)

var (
	n1  *NetworkObject
	g1  *NetworkObjectGroup
	ftd *FTD
)

func setupTestAccessRuleObjects(t *testing.T) error {
	var err error

	ftd, err = initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return err
	}

	tearDownTestAccessRuleObjects(t)

	n1 = new(NetworkObject)
	n1.Name = "testObj001"
	n1.SubType = "HOST"
	n1.Value = "1.1.1.1"

	err = ftd.CreateNetworkObject(n1, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return err
	}

	g1 = new(NetworkObjectGroup)
	g1.Name = "testObjGroup001"
	g1.Objects = append(g1.Objects, n1.Reference())

	err = ftd.CreateNetworkObjectGroup(g1, DuplicateActionReplace)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return err
	}

	return nil
}

func tearDownTestAccessRuleObjects(t *testing.T) error {
	var err error

	if g1 != nil {
		err = ftd.DeleteNetworkObjectGroup(g1)
		if err != nil {
			t.Errorf("error: %s\n", err)
			return err
		}
	}

	if n1 != nil {
		err = ftd.DeleteNetworkObject(n1)
		if err != nil {
			t.Errorf("error: %s\n", err)
			return err
		}
	}
	return nil
}

func TestGetAccessRules(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	_, err = ftd.GetAccessRules("default", 0)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}
}

func TestCreateAccessRules(t *testing.T) {
	var err error

	ftd, err := initTest()
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}
	a := new(AccessRule)
	a.Name = "testPolicy001"
	a.RuleAction = RuleActionPermit
	a.EventLogAction = LogActionNone

	err = ftd.CreateAccessRule(a, "default")
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	if a.ID == "" {
		t.Errorf("ID was not set\n")
	}

	err = ftd.DeleteAccessRule(a)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}
}

func TestCombinedCreateAccessRules(t *testing.T) {
	var err error

	err = setupTestAccessRuleObjects(t)
	if err != nil {
		t.Errorf("error: %s\n", err)
		return
	}

	a := new(AccessRule)
	a.Name = "testPolicy001"
	a.RuleAction = RuleActionPermit
	a.EventLogAction = LogActionNone
	a.DestinationNetworks = append(a.DestinationNetworks, n1.Reference())
	a.DestinationNetworks = append(a.DestinationNetworks, g1.Reference())

	err = ftd.CreateAccessRule(a, "default")
	if err != nil {
		t.Errorf("error: %s\n", err)
		tearDownTestAccessRuleObjects(t)
		return
	}

	//spew.Dump(a)

	err = ftd.DeleteAccessRule(a)
	if err != nil {
		t.Errorf("error: %s\n", err)
		tearDownTestAccessRuleObjects(t)
		return
	}

	tearDownTestAccessRuleObjects(t)

}
