package prismacloud_client

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestMakeRequestHttpClientIsSet(t *testing.T) {
	client, close := makeTestClient(t, "", "", "", 200, "")
	defer close()

	assert.NotNil(t, client.http)
}

func TestMakeRequestRefreshesAuth(t *testing.T) {
	verb := "not-used"
	path := "/not-used"

	client, close := makeTestClient(t, verb, path, "", 200, "")
	client.refreshingAuth = true
	defer close()

	called := false
	client.refreshAuth = func() { called = true }

	client.makeRequest(verb, path, nil, nil)

	assert.True(t, called)
}

func TestMakeRequestGet(t *testing.T) {
	verb := "GET"
	path := "/my/special/path"
	response := `{"Message": "ok"}`

	client, close := makeTestClient(t, verb, path, "", 200, response)
	client.refreshingAuth = true
	defer close()

	retval := new(TestMakeRequestBody)
	client.makeRequest(verb, path, nil, retval)

	assert.Equal(t, "ok", retval.Message)
}

func TestMakeRequestPost(t *testing.T) {
	path := "/my/special/path"
	verb := "POST"
	message := `{"Message": "ok"}`

	client, close := makeTestClient(t, verb, path, message, 200, "")
	client.refreshingAuth = true
	defer close()

	req := bytes.NewBuffer([]byte(message))
	client.refreshingAuth = true
	client.makeRequest(verb, path, req, nil)
}

func TestMakeRequestErrorsOnInvalidVerb(t *testing.T) {
	verb := "  "
	path := "/not-used"

	client, close := makeTestClient(t, verb, path, "", 200, "")
	client.refreshingAuth = true
	defer close()

	err := client.makeRequest(verb, path, nil, nil)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failure creating request")
}

func TestMakeRequestErrorsOnTimeout(t *testing.T) {
	verb := "GET"
	path := "/not-used"

	client, close := makeTestClient(t, verb, path, "", 200, "")
	client.refreshingAuth = true
	close() //close server so we get connection refused

	err := client.makeRequest(verb, path, nil, nil)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Failure making request")
}

func TestMakeRequestErrorsBadResponse(t *testing.T) {
	verb := "GET"
	path := "/not-used"

	client, close := makeTestClient(t, verb, path, "", 401, "")
	client.refreshingAuth = true
	defer close()

	client.readResponse = func(body io.ReadCloser) ([]byte, error) {
		return nil, fmt.Errorf("Injected error from TestMakeRequestErrorsBadResponse")
	}

	err := client.makeRequest(verb, path, nil, nil)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error reading response")
}

func TestMakeRequestErrorsBadJson(t *testing.T) {
	verb := "GET"
	path := "/not-used"

	client, close := makeTestClient(t, verb, path, "", 200, `{ "incomplete": "json}`)
	client.refreshingAuth = true
	defer close()

	retval := new(TestMakeRequestBody)
	err := client.makeRequest(verb, path, nil, retval)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error decoding json")
}

func TestMakeRequestNone200ReturnsError(t *testing.T) {
	verb := "GET"
	path := "/not-used"

	client, close := makeTestClient(t, verb, path, "", 400, "")
	client.refreshingAuth = true
	defer close()

	err := client.makeRequest(verb, path, nil, nil)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error: injected_test_error")
}

func TestMakeRequestEmptyResponse(t *testing.T) {
	verb := "GET"
	path := "/not-used"

	client, close := makeTestClient(t, verb, path, "", 200, "")
	client.refreshingAuth = true
	defer close()

	err := client.makeRequest(verb, path, nil, nil)

	assert.Nil(t, err)
}
