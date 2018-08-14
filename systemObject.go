package goftd

import (
	"fmt"

	"github.com/golang/glog"
)

// GetNetworkAny Returns the 0.0.0.0/0 object
func (f *FTD) GetNetworkAny() (*NetworkObject, error) {
	obj, err := f.getNetworkObjectBy("name:0.0.0.0", 1)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, err
	}
	if len(obj) != 1 {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return nil, fmt.Errorf("Error: %s", err)
	}

	return obj[0], nil
}
