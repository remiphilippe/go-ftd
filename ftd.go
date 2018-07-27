package goftd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/golang/glog"

	"github.com/davecgh/go-spew/spew"
)

type passwordGrant struct {
	Username string
	Password string
}

//TODO: implement this?
type customGrant struct {
}

// FTD struct holding the FTD object
type FTD struct {
	// Hostname or IP address
	Hostname string
	// Define authorization type as password or custom
	GrantType string
	// store access token and refresh token
	accessToken  string
	refreshToken string
	expiresAt    time.Time

	passwordGrant *passwordGrant
	customGrant   *customGrant

	debug bool
}

// Links Embedded links
type Links struct {
	Self string `json:"self,omitempty"`
}

// ReferenceObject FTD reference object
type ReferenceObject struct {
	ID      string `json:"id,omitempty"`
	Version string `json:"version,omitempty"`
	Name    string `json:"name"`
	Type    string `json:"type"`
}

func (f *FTD) updateToken() error {
	if f.GrantType == "" {
		return fmt.Errorf("grant is not correctly initialized")
	}

	if f.GrantType == grantTypePassword && f.passwordGrant == nil {
		return fmt.Errorf("grant is not correctly initialized")
	} else if f.GrantType == grantTypeCustom && f.customGrant == nil {
		return fmt.Errorf("grant is not correctly initialized")
	}

	req := make(map[string]string)
	var res map[string]interface{}

	if f.GrantType == grantTypePassword {

		req["grant_type"] = grantTypePassword
		req["username"] = f.passwordGrant.Username
		req["password"] = f.passwordGrant.Password
	}

	data, err := f.Post("fdm/token", req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		if f.debug {
			glog.Errorf("Error: %s\n", err)
		}
		return err
	}

	if _, ok := res["access_token"]; ok {
		f.accessToken = res["access_token"].(string)
	} else {
		return fmt.Errorf("missing access_token in reply")
	}

	if _, ok := res["refresh_token"]; ok {
		f.refreshToken = res["refresh_token"].(string)
	} else {
		return fmt.Errorf("missing refresh_token in reply")
	}

	if _, ok := res["expires_in"]; ok {
		f.expiresAt = time.Now().Add(time.Second * time.Duration(res["expires_in"].(float64)))
	} else {
		return fmt.Errorf("missing expires_in in reply")
	}

	return nil
}

// NewFTD returns an initilized FTD struct
func NewFTD(hostname string, param map[string]string) (*FTD, error) {
	f := new(FTD)
	f.Hostname = hostname
	f.debug = false

	if _, ok := param["debug"]; ok {
		if param["debug"] == "true" {
			f.debug = true
		}
	}

	if _, ok := param["grant_type"]; ok {
		if param["grant_type"] == grantTypePassword || param["grant_type"] == grantTypeCustom {
			f.GrantType = param["grant_type"]
			if param["grant_type"] == grantTypePassword {
				f.passwordGrant = new(passwordGrant)
				if _, ok := param["username"]; ok {
					f.passwordGrant.Username = param["username"]
				} else {
					if f.debug {
						glog.Errorf("username is mandatory for grant type = %s\n", grantTypePassword)
					}
					return nil, fmt.Errorf("username is mandatory for grant type = %s", grantTypePassword)
				}

				if _, ok := param["password"]; ok {
					f.passwordGrant.Password = param["password"]
				} else {
					if f.debug {
						glog.Errorf("password is mandatory for grant type = %s\n", grantTypePassword)
					}
					return nil, fmt.Errorf("password is mandatory for grant type = %s", grantTypePassword)
				}
			}
		} else {
			if f.debug {
				glog.Errorf("unknown grant type: %s\n", param["grant_type"])
			}
			return nil, fmt.Errorf("unknown grant type: %s", param["grant_type"])
		}
	}

	err := f.updateToken()
	if err != nil {
		return nil, err
	}

	return f, nil
}

// Post POST to ASA API
func (f *FTD) Post(endpoint string, ftdReq interface{}) (bodyText []byte, err error) {
	var authenticating bool

	if endpoint != "fdm/token" {
		authenticating = false
	} else {
		authenticating = true
	}

	uri := url.URL{
		Host:   f.Hostname,
		Scheme: "https",
		Path:   "api/fdm/v1/" + endpoint,
	}

	var jsonReq []byte
	var body io.Reader

	if ftdReq != nil {
		jsonReq, err = json.Marshal(ftdReq)
		if err != nil {
			glog.Errorf("POST - marshall error: %s\n", err)
			return nil, err
		}
		body = bytes.NewBuffer(jsonReq)
	} else {
		body = nil
	}

	//spew.Dump(string(jsonReq))

	req, err := http.NewRequest("POST", uri.String(), body)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	if !authenticating {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", f.accessToken))
	}

	resp, err := client.Do(req)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	bodyText, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("POST - readall error: %s\n", err)
		spew.Dump(resp)
		return nil, err
	}

	glog.Infof("Response: %s\n", strconv.Itoa(resp.StatusCode))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = parseResponse(bodyText, authenticating)
		if err != nil {
			if f.debug {
				glog.Errorf("POST - parse response error: %s\n", err)
			}
			return nil, err
		}

		return nil, fmt.Errorf("response code: %d", resp.StatusCode)
	}

	return bodyText, nil
}

// Put PUT to ASA API
func (f *FTD) Put(endpoint string, ftdReq interface{}) (bodyText []byte, err error) {
	uri := url.URL{
		Host:   f.Hostname,
		Scheme: "https",
		Path:   "api/fdm/v1/" + endpoint,
	}

	var jsonReq []byte
	var body io.Reader

	if ftdReq != nil {
		jsonReq, err = json.Marshal(ftdReq)
		if err != nil {
			glog.Errorf("POST - marshall error: %s\n", err)
			return nil, err
		}
		body = bytes.NewBuffer(jsonReq)
	} else {
		body = nil
	}

	//spew.Dump(string(jsonReq))

	req, err := http.NewRequest("PUT", uri.String(), body)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if endpoint != "fdm/token" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", f.accessToken))
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	bodyText, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("POST - readall error: %s\n", err)
		spew.Dump(resp)
		return nil, err
	}

	glog.Infof("Response: %s\n", strconv.Itoa(resp.StatusCode))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = parseResponse(bodyText, false)
		if err != nil {
			if f.debug {
				glog.Errorf("POST - parse response error: %s\n", err)
			}
			return nil, err
		}

		return nil, fmt.Errorf("response code: %d", resp.StatusCode)
	}

	return bodyText, nil
}

// Get GET to ASA API
func (f *FTD) Get(endpoint string) (bodyText []byte, err error) {
	uri := url.URL{
		Host:   f.Hostname,
		Scheme: "https",
		Path:   "api/fdm/v1/" + endpoint,
	}

	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if f.accessToken == "" {
		return nil, fmt.Errorf("accessToken is not set, did you initialize correctly?")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", f.accessToken))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	bodyText, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("GET - readall error: %s\n", err)
		spew.Dump(resp)
		return nil, err
	}

	log.Print("Response: " + strconv.Itoa(resp.StatusCode))
	//spew.Dump(string(bodyText))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = parseResponse(bodyText, false)
		if err != nil {
			if f.debug {
				glog.Errorf("GET - parse response error: %s\n", err)
			}
			return nil, err
		}
	}

	return bodyText, nil
}

// Delete DELETE to ASA API
func (f *FTD) Delete(endpoint string) (err error) {
	uri := url.URL{
		Host:   f.Hostname,
		Scheme: "https",
		Path:   "api/fdm/v1/" + endpoint,
	}

	req, err := http.NewRequest("DELETE", uri.String(), nil)
	if err != nil {
		log.Print(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if f.accessToken == "" {
		return fmt.Errorf("accessToken is not set, did you initialize correctly?")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", f.accessToken))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("DELETE - readall error: %s\n", err)
		spew.Dump(resp)
		return err
	}

	log.Print("Response: " + strconv.Itoa(resp.StatusCode))
	err = parseResponse(bodyText, false)
	if err != nil {
		if f.debug {
			glog.Errorf("DELETE - parse response error: %s\n", err)
		}
		return err
	}

	return nil
}
