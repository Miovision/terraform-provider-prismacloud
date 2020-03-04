package prismacloud_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type PrismaCloudClient struct {
	baseUrl        string
	accessKey      string
	secretKey      string
	http           *http.Client
	refreshingAuth bool
	refreshAuth    func()
	makeRequest    func(string, string, io.Reader, interface{}) error
	readResponse   func(io.ReadCloser) ([]byte, error)
	token          string
}

func MakePrismaCloudClient(baseUrl string, accessKey string, secretKey string) *PrismaCloudClient {
	c := new(PrismaCloudClient)
	c.baseUrl = baseUrl
	c.accessKey = accessKey
	c.secretKey = secretKey
	c.http = &http.Client{Timeout: 60 * time.Second}
	c.refreshingAuth = false

	c.refreshAuth = func() {
		if c.refreshingAuth {
			return
		}
		c.refreshingAuth = true

		// TODO Need to check if token is still valid and refresh if it is not
		if c.token == "" {
			resp := new(AuthResponse)
			c.Post("/login", &AuthRequest{c.accessKey, c.secretKey}, resp)
			c.token = resp.Token
		}
		c.refreshingAuth = false
	}

	c.makeRequest = func(verb string, path string, body io.Reader, retval interface{}) error {
		c.refreshAuth()
		req, err := http.NewRequest(verb, c.getUrl(path), body)
		if err != nil {
			return fmt.Errorf("Failure creating request for method '%s' and path '%s'. %+v", verb, path, err)
		}

		if c.token != "" {
			println("ADDING AUTH HEADER!")
			req.Header.Add("x-redlock-auth", c.token)
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("accept", "application/json; charset=UTF-8")
		resp, err := c.http.Do(req)
		if err != nil {
			return fmt.Errorf("Failure making request for method '%s' and path '%s'. %+v", verb, path, err)
		}
		defer resp.Body.Close()

		fmt.Printf("request: %+v\n-----\n", req)
		fmt.Printf("response: %+v\n-----\n", resp)

		byteBody, err := c.readResponse(resp.Body)
		if err != nil {
			return fmt.Errorf("Error reading response. %+v %+v", resp, err)
		}

		switch resp.StatusCode {
		case http.StatusOK:
			break
		case http.StatusUnauthorized:
			return fmt.Errorf("You are not authorized.  Please check your API credentials. Response was: %+v, Body: %s.", resp, string(byteBody))
		default:
			status := []PrismaCloudError{}
			err = json.NewDecoder(strings.NewReader(resp.Header.Get("X-Redlock-Status"))).Decode(&status)
			if err != nil {
				return fmt.Errorf("Fatal error decoding error response json. Request: +%v, Response: %+v. Error: %+v. Body: %s", req, resp, err, string(byteBody))
			}
			return fmt.Errorf("Error: %s.\nResponse was: %+v.\nBody: %s", status[0].Reason, resp, string(byteBody))
		}

		fmt.Printf("response body: %s\n-----\n", string(byteBody))
		if retval != nil {
			err = json.NewDecoder(bytes.NewBuffer(byteBody)).Decode(retval)
			if err != nil {
				return fmt.Errorf("Error decoding json. Response: %+v. Error: %+v. Body: %s", resp, err, string(byteBody))
			}
			fmt.Printf("response object: %+v\n------------------------------\n", retval)
		}
		return nil
	}

	c.readResponse = func(body io.ReadCloser) ([]byte, error) {
		return ioutil.ReadAll(body)
	}

	return c
}

func (c PrismaCloudClient) getUrl(path string) string {
	return c.baseUrl + path
}

func (c PrismaCloudClient) Get(path string, retval interface{}) error {
	return c.makeRequest("GET", path, nil, retval)
}

func (c PrismaCloudClient) Post(path string, body interface{}, resp interface{}) error {
	requestBody, _ := json.Marshal(body)
	print("POST request body: ")
	println(string(requestBody))
	println("------------------------------")
	return c.makeRequest("POST", path, bytes.NewBuffer(requestBody), resp)
}

func (c PrismaCloudClient) Put(path string, body interface{}) error {
	requestBody, _ := json.Marshal(body)
	print("PUT request body: ")
	println(string(requestBody))
	println("------------------------------")
	return c.makeRequest("PUT", path, bytes.NewBuffer(requestBody), nil)
}

func (c PrismaCloudClient) Delete(path string) error {
	return c.makeRequest("DELETE", path, nil, nil)
}
