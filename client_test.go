package gitlab_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kryabinin/go-gitlab"
)

func TestClient_SendRequest(t *testing.T) {
	t.Run("positive case", func(t *testing.T) {
		var (
			token       = "test_token"
			baseUrl     = "http://gitlab.test.com/api/v4"
			path        = "/test/path"
			method      = http.MethodGet
			expResponse = []byte(`{"test": "passed"}`)
		)

		httpClient := new(gitlab.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
			req, ok := args.Get(0).(*http.Request)
			assert.True(t, ok)

			assert.Equal(t, method, req.Method)
			assert.Equal(t, baseUrl+"/"+path, req.URL.String())

			assert.Equal(t, token, req.Header.Get("Private-Token"))
			assert.Equal(t, "application/json; charset=utf-8", req.Header.Get("Content-Type"))
		}).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader(expResponse)),
			StatusCode: http.StatusOK,
		}, nil)

		client := gitlab.NewClient(
			token,
			gitlab.WithBaseUrl(baseUrl),
			gitlab.WithHttpClient(httpClient),
		)

		resp, err := client.SendRequest(context.Background(), method, path, nil)
		assert.NoError(t, err)
		assert.Equal(t, expResponse, resp)
	})

	t.Run("error on sending http request", func(t *testing.T) {
		expErr := errors.New("test error")

		httpClient := new(gitlab.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, expErr)

		client := gitlab.NewClient(
			"test_token",
			gitlab.WithHttpClient(httpClient),
		)

		resp, err := client.SendRequest(context.Background(), http.MethodGet, "/test/path", nil)
		assert.Equal(t, []byte(nil), resp)
		assert.True(t, errors.Is(err, expErr))
	})

	t.Run("non 200 status", func(t *testing.T) {
		expStatus := http.StatusInternalServerError

		httpClient := new(gitlab.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
			StatusCode: expStatus,
		}, nil)

		client := gitlab.NewClient(
			"test_token",
			gitlab.WithHttpClient(httpClient),
		)

		resp, err := client.SendRequest(context.Background(), http.MethodGet, "/test/path", nil)
		assert.Equal(t, []byte(nil), resp)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("gitlab respond with %d status code", expStatus))
	})

	t.Run("error on read response body", func(t *testing.T) {
		expStatus := http.StatusInternalServerError

		httpClient := new(gitlab.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
			StatusCode: expStatus,
		}, nil)

		client := gitlab.NewClient(
			"test_token",
			gitlab.WithHttpClient(httpClient),
		)

		resp, err := client.SendRequest(context.Background(), http.MethodGet, "/test/path", nil)
		assert.Equal(t, []byte(nil), resp)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("gitlab respond with %d status code", expStatus))
	})
}
