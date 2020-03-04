package prismacloud_client

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

//Sets up the PrismaCloudClient with a mock makeRequest and refreshAuth function.
//Use to isolate your tests from these two core functions and focus on how you call them.
func makeSimpleMockTestClient(t *testing.T, response interface{}) (*PrismaCloudClient, *MakeRequestMock) {
	client, mockRequest := makeMockTestClient(t, response)
	client.refreshAuth = func() {}

	matchAny := mock.MatchedBy(func(_ interface{}) bool { return true })

	mockRequest.
		On("makeRequest", mock.AnythingOfType("string"), mock.AnythingOfType("string"), matchAny, matchAny).
		Once()

	return client, mockRequest
}

func makeSimpleErrorTestClient(t *testing.T, err error) (*PrismaCloudClient, *MakeRequestMock) {
	client, mockRequest := makeSimpleMockTestClient(t, nil)
	mockRequest.Error = err
	return client, mockRequest
}

func TestGetForwardsToMakeRequest(t *testing.T) {
	path := "/my/cool/path"
	expected := TestMakeRequestBody{"tester"}

	client, mockRequest := makeSimpleMockTestClient(t, expected)

	actual := new(TestMakeRequestBody)
	client.Get(path, actual)

	assert.Equal(t, &expected, actual)
	assert.Equal(t, 1, len(mockRequest.Calls))
	assert.True(t, mockRequest.AssertCalled(t, "makeRequest", "GET", path, nil, actual))
}

func TestGetReturnsError(t *testing.T) {
	expectedError := fmt.Errorf("Injecting error for TestGetReturnsError")

	client, _ := makeSimpleErrorTestClient(t, expectedError)

	actualError := client.Get("", nil)

	assert.Equal(t, expectedError, actualError)
}

func TestPostForwardsToMakeRequest(t *testing.T) {
	path := "/my/cool/path"

	client, mockRequest := makeSimpleMockTestClient(t, nil)

	request := TestMakeRequestBody{"My Post Message"}
	requestJson := fmt.Sprintf(`{"Message":"%s"}`, request.Message)
	client.Post(path, request, nil)

	assert.Equal(t, 1, len(mockRequest.Calls))
	assert.True(t, mockRequest.AssertCalled(t, "makeRequest", "POST", path, bytes.NewBuffer([]byte(requestJson)), nil))
}

func TestPostReturnsError(t *testing.T) {
	expectedError := fmt.Errorf("Injecting error for TestPostReturnsError")

	client, _ := makeSimpleErrorTestClient(t, expectedError)

	actualError := client.Post("", nil, nil)

	assert.Equal(t, expectedError, actualError)
}

func TestPutForwardsToMakeRequest(t *testing.T) {
	path := "/my/cool/path"

	client, mockRequest := makeSimpleMockTestClient(t, nil)

	request := TestMakeRequestBody{"My Post Message"}
	requestJson := fmt.Sprintf(`{"Message":"%s"}`, request.Message)
	client.Put(path, request)

	assert.Equal(t, 1, len(mockRequest.Calls))
	assert.True(t, mockRequest.AssertCalled(t, "makeRequest", "PUT", path, bytes.NewBuffer([]byte(requestJson)), nil))
}

func TestPutReturnsError(t *testing.T) {
	expectedError := fmt.Errorf("Injecting error for TestPutReturnsError")

	client, _ := makeSimpleErrorTestClient(t, expectedError)

	actualError := client.Put("", nil)

	assert.Equal(t, expectedError, actualError)
}

func TestDeleteForwardsToMakeRequest(t *testing.T) {
	path := "/my/cool/path"

	client, mockRequest := makeSimpleMockTestClient(t, nil)

	client.Delete(path)

	assert.Equal(t, 1, len(mockRequest.Calls))
	assert.True(t, mockRequest.AssertCalled(t, "makeRequest", "DELETE", path, nil, nil))
}

func TestDeleteReturnsError(t *testing.T) {
	expectedError := fmt.Errorf("Injecting error for TestDeleteReturnsError")

	client, _ := makeSimpleErrorTestClient(t, expectedError)

	actualError := client.Delete("")

	assert.Equal(t, expectedError, actualError)
}
