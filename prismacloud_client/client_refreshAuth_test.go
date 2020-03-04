package prismacloud_client

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRefreshAuthDoesNothingWhenAlreadyRefreshing(t *testing.T) {
	client, mock := makeMockTestClient(t, nil)

	client.refreshingAuth = true

	client.refreshAuth()

	assert.True(t, mock.AssertNotCalled(t, "makeRequest"))
}

func TestRefreshAuthDoesNothingWhenTokenAlreadySet(t *testing.T) {
	client, mock := makeMockTestClient(t, nil)
	client.token = "jwt-token"

	client.refreshAuth()

	assert.True(t, mock.AssertNotCalled(t, "makeRequest"))
}

func TestRefreshAuthSetsRefreshingFlag(t *testing.T) {
	client, mockRequest := makeMockTestClient(t, nil)

	//We know that if we need to refresh we should only call makeRequest once.
	//Relies on MakeRequestMock.makeRequest calling refreshAuth like PrismaClient.refreshAuth does.
	mockRequest.
		On("makeRequest", "POST", "/login", mock.AnythingOfType("*bytes.Buffer"), mock.AnythingOfType("*prismacloud_client.AuthResponse")).
		Return(nil).
		Once()

	assert.False(t, client.refreshingAuth)
	client.refreshAuth()
	assert.False(t, client.refreshingAuth)
}

func TestRefreshAuthSetsToken(t *testing.T) {
	response := AuthResponse{"my-jwt-token", "Success!", []Customer{}}
	client, mockRequest := makeMockTestClient(t, response)

	mockRequest.
		On("makeRequest", "POST", "/login", mock.AnythingOfType("*bytes.Buffer"), mock.AnythingOfType("*prismacloud_client.AuthResponse")).
		Return(nil).
		Once()

	client.refreshAuth()

	assert.Equal(t, client.token, response.Token)
}
