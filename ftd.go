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

	Insecure bool

	// store access token and refresh token
	accessToken  string
	refreshToken string
	expiresAt    time.Time

	passwordGrant *passwordGrant
	customGrant   *customGrant

	debug bool
}

type requestParameters struct {
	// Request for POST / PUT
	FTDRequest interface{}
	// URI Query if needed (GET)
	URIQuery map[string]string
	// Paging parameters for GET
	PageStart int
	PageLimit int
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

// Paging Paging Information
type Paging struct {
	Prev   []string `json:"prev,omitempty"`
	Next   []string `json:"next,omitempty"`
	Limit  int      `json:"limit,omitempty"`
	Offset int      `json:"offset,omitempty"`
	Count  int      `json:"count,omitempty"`
	Pages  int      `json:"pages,omitempty"`
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

	data, err := f.Post(apiTokenEndpoint, req)
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
	f.Insecure = false

	if _, ok := param["debug"]; ok {
		if param["debug"] == "true" {
			f.debug = true
		}
	}

	if _, ok := param["insecure"]; ok {
		if param["insecure"] == "true" {
			f.Insecure = true
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

func (f *FTD) request(endpoint, method string, r *requestParameters) (bodyText []byte, err error) {
	var authenticating bool
	var req *http.Request
	var jsonReq []byte
	var body io.Reader

	if endpoint == apiTokenEndpoint && method == "POST" {
		authenticating = true
	} else {
		authenticating = false
	}

	uri := url.URL{
		Host:   f.Hostname,
		Scheme: "https",
		Path:   apiBasePath + endpoint,
	}

	switch method {
	case apiPOST, apiPUT:
		if r != nil && r.FTDRequest != nil {
			jsonReq, err = json.Marshal(r.FTDRequest)
			if err != nil {
				glog.Errorf("request - marshall error: %s\n", err)
				return nil, err
			}
			body = bytes.NewBuffer(jsonReq)
		} else {
			body = nil
		}

		req, err = http.NewRequest(method, uri.String(), body)
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

	case apiGET, apiDELETE:
		req, err = http.NewRequest(method, uri.String(), nil)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		q := req.URL.Query()
		if r != nil {
			for k, v := range r.URIQuery {
				q.Add(k, v)
			}
		}
		// if method == apiGET {
		// 	if r != nil && r.PageLimit > 0 {
		// 		q.Add("limit", string(r.PageLimit))
		// 	}

		// 	if r != nil && r.PageStart > 0 {
		// 		q.Add("start", string(r.PageStart))
		// 	}
		// }

		req.URL.RawQuery = q.Encode()

	default:
		return nil, fmt.Errorf("Unknown Method %s", method)
	}

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
		glog.Errorf("request - readall error: %s\n", err)
		spew.Dump(resp)
		return nil, err
	}

	glog.Infof("Response: %s\n", strconv.Itoa(resp.StatusCode))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = parseResponse(bodyText, authenticating)
		if err != nil {
			// if f.debug {
			// 	glog.Errorf("POST - parse response error: %s\n", err)
			// }
			return nil, err
		}

		return nil, fmt.Errorf("response code: %d", resp.StatusCode)
	}

	return bodyText, nil
}

// Post POST to ASA API
func (f *FTD) Post(endpoint string, ftdReq interface{}) (bodyText []byte, err error) {
	r := requestParameters{
		FTDRequest: ftdReq,
	}
	return f.request(endpoint, apiPOST, &r)
}

// Put PUT to ASA API
func (f *FTD) Put(endpoint string, ftdReq interface{}) (bodyText []byte, err error) {
	r := requestParameters{
		FTDRequest: ftdReq,
	}
	return f.request(endpoint, apiPUT, &r)
}

// Get GET to ASA API
func (f *FTD) Get(endpoint string, uriQuery map[string]string) (bodyText []byte, err error) {
	r := requestParameters{
		URIQuery: uriQuery,
	}
	return f.request(endpoint, apiGET, &r)
}

// Delete DELETE to ASA API
func (f *FTD) Delete(endpoint string) (err error) {
	_, err = f.request(endpoint, apiDELETE, nil)
	return err
}
