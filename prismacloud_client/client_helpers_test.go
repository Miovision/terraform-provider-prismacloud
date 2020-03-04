package prismacloud_client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TestMakeRequestBody struct {
	Message string
}

//Sets up the PrismaCloudClient with a fake HTTP server.  This function will check that the correct body was sent
//and will return the `response` to the client.
func makeTestClient(t *testing.T, verb string, path string, body string, statusCode int, response string) (*PrismaCloudClient, func()) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, verb, req.Method)
		assert.Equal(t, path, req.URL.String())

		reqBytes, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, body, string(reqBytes))

		rw.Header().Set("X-Redlock-Status", `[{"i18nKey":"injected_test_error","severity":"error","subject":null}]`)
		rw.WriteHeader(statusCode)
		rw.Write([]byte(response))
	}))

	client := MakePrismaCloudClient(server.URL, "fake", "fake")
	client.http = server.Client()

	return client, server.Close
}

type MakeRequestMock struct {
	mock.Mock
	client      *PrismaCloudClient
	ReturnValue interface{}
	Error       error
}

func (m *MakeRequestMock) makeRequest(verb string, path string, body io.Reader, retval interface{}) error {
	m.Called(verb, path, body, retval)
	m.client.refreshAuth() //Important to call this to mimic what the real function does

	if m.ReturnValue != nil {
		reflRetval := reflect.ValueOf(retval)
		ptrRetval := reflect.Indirect(reflRetval)

		fmt.Printf("ReflRetval: %+v %v, PtrRetval: %v %v\n", reflRetval, reflRetval.CanSet(), ptrRetval, ptrRetval.CanSet())
		ptrRetval.Set(reflect.ValueOf(m.ReturnValue))
	}
	return m.Error
}

//Sets up the PrismaCloudClient with a mock makeRequest function.
func makeMockTestClient(t *testing.T, response interface{}) (*PrismaCloudClient, *MakeRequestMock) {
	client := MakePrismaCloudClient("fake", "fake", "fake")

	reqMock := new(MakeRequestMock)
	reqMock.client = client
	reqMock.ReturnValue = response
	client.makeRequest = reqMock.makeRequest

	return client, reqMock
}
