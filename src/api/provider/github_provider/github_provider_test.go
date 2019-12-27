package github_provider

import (
	"errors"
	"fmt"
	"github.com/Sisi55/GoConsumeGitHubAPI/src/api/clients/restclient"
	"github.com/Sisi55/GoConsumeGitHubAPI/src/api/config"
	"github.com/Sisi55/GoConsumeGitHubAPI/src/api/domain/github"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "Authorization", headerAuthorization)
	assert.EqualValues(t, "token %s", headerAuthorizationFormat)
	assert.EqualValues(t, "https://api.github.com/user/repos", urlCreateRepo)
}

func TestGetAuthorizationHeader(t *testing.T) {
	header := getAuthorizationHeader("abc123")
	assert.EqualValues(t, "token abc123", header)
}

func TestDefer(t *testing.T) {
	defer fmt.Println("1")
	defer fmt.Println("2")
	defer fmt.Println("3")

	fmt.Println("function's body")
}

func TestCreateRepoErrorRestClient(t *testing.T) {
	restclient.FlushMockups()
	//restclient.StartMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response:   nil,
		Err:        errors.New("invalid restclient response"),
	})

	response, err := CreateRepo(config.GetGithubAccessToken(), github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid restclient response", err.Message)

	//restclient.StopMockups()
	//
	//response, err = CreateRepo("", github.CreateRepoRequest{})
	//assert.Nil(t, response)
	//assert.NotNil(t, err)
}

func TestCreateRepoInvalidResponseBody(t *testing.T) {
	restclient.FlushMockups()
	//restclient.StartMockups()
	invaliCloser, _ := os.Open("-asf3")
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       invaliCloser,
		},
		//Err:        errors.New("invalid restclient response"),
	})

	response, err := CreateRepo(config.GetGithubAccessToken(), github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid response body", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
}

func TestCreateRepoInvalidErrorInterface(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message":1}`)),
		},
		//Err:        errors.New("invalid restclient response"),
	})

	response, err := CreateRepo(config.GetGithubAccessToken(), github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid json response body", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
}

func TestCreateRepoUnauthorized(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Bad credentials","documentation_url": "https://developer.github.com/v3"}`)),
		},
		//Err:        errors.New("invalid restclient response"),
	})

	response, err := CreateRepo(config.GetGithubAccessToken(), github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Bad credentials", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
}

func TestCreateRepoSuccessInvalidResponse(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": "abc123"}`)),
		},
		//Err:        errors.New("invalid restclient response"),
	})

	response, err := CreateRepo(config.GetGithubAccessToken(), github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Bad credentials", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
}

func TestCreateRepoNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123,"name":"...","full_name":"..."}`)),
		},
		//Err:        errors.New("invalid restclient response"),
	})

	response, err := CreateRepo(config.GetGithubAccessToken(), github.CreateRepoRequest{})
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, 123, response.Id)
	assert.EqualValues(t, "...", response.Name)
	assert.EqualValues(t, "...", response.FullName)

}
