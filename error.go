package goftd

import (
	"encoding/json"
	"fmt"
)

// FTDMessage  Error message returned by API
type FTDMessage struct {
	Description string
	Code        string
	Location    string
}

// FTDError Error returned by API
type FTDError struct {
	Severity string       `json:"severity"`
	Key      string       `json:"key"`
	Message  []FTDMessage `json:"messages"`
}

func (fe FTDError) Error() string {
	return fmt.Sprintf("%s: %s with messages %+v", fe.Severity, fe.Key, fe.Message)
}

func parseResponse(bodyText []byte, authenticating bool) (err error) {
	//spew.Dump(string(bodyText))
	if len(bodyText) > 0 {
		if !authenticating {
			var v struct {
				Error *FTDError `json:"error"`
			}

			//log.Print("Response: " + string(bodyText))

			err = json.Unmarshal(bodyText, &v)
			if err != nil {
				return err
			}

			return v.Error
		}

		var v map[string]interface{}

		//log.Print("Response: " + string(bodyText))

		err = json.Unmarshal(bodyText, &v)
		if err != nil {
			return err
		}

		return fmt.Errorf("error getting token: %s", v["message"].(string))

	}
	return nil
}
