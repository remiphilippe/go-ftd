package goftd

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
)

func initTest() (*FTD, error) {
	params := make(map[string]string)
	params["grant_type"] = "password"
	params["username"] = os.Getenv("FTD_USER")
	params["password"] = os.Getenv("FTD_PASSWORD")
	params["debug"] = "true"

	ftd, err := NewFTD(os.Getenv("FTD_HOST"), params)
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return nil, err
	}

	return ftd, nil
}

func TestToken(t *testing.T) {
	ftd, err := initTest()
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return
	}

	spew.Dump(ftd)
}
