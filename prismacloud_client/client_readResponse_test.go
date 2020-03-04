package prismacloud_client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

type TestErrorReader struct {
	io.Reader
	Err error
}

func (r TestErrorReader) Read(p []byte) (n int, err error) {
	return 0, r.Err
}

func TestReadResponseThrowsError(t *testing.T) {
	client, _ := makeMockTestClient(t, nil)
	expectedError := fmt.Errorf("Injecting test error for TestReadResponseThrowsError")
	resp := ioutil.NopCloser(TestErrorReader{Err: expectedError})

	_, err := client.readResponse(resp)

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
}

func TestReadResponseReadsFullResponse(t *testing.T) {
	client, _ := makeMockTestClient(t, nil)
	message := "My full testing message"
	resp := ioutil.NopCloser(strings.NewReader(message))

	bytes, err := client.readResponse(resp)

	assert.Nil(t, err)
	assert.Equal(t, message, string(bytes))
}
