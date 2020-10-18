package gitlab_test

import (
	"bytes"
	"context"
	"encoding/json"
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
			gitlab.ClientOptionBaseUrl(baseUrl),
			gitlab.ClientOptionHttpClient(httpClient),
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
			gitlab.ClientOptionHttpClient(httpClient),
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
			gitlab.ClientOptionHttpClient(httpClient),
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
			gitlab.ClientOptionHttpClient(httpClient),
		)

		resp, err := client.SendRequest(context.Background(), http.MethodGet, "/test/path", nil)
		assert.Equal(t, []byte(nil), resp)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("gitlab respond with %d status code", expStatus))
	})
}

func TestClient_GetUserByID(t *testing.T) {
	t.Run("positive case", func(t *testing.T) {
		var (
			userID  = 5
			baseUrl = "http://gitlab.test.com/api/v4"
		)

		httpClient := new(gitlab.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
			req, ok := args.Get(0).(*http.Request)

			assert.True(t, ok)
			assert.Equal(t, baseUrl+"/"+fmt.Sprintf("/users/%d", userID), req.URL.String())
		}).Return(&http.Response{
			Body: ioutil.NopCloser(
				bytes.NewReader(
					[]byte(fmt.Sprintf("{\"id\": %d}", userID)),
				),
			),
			StatusCode: http.StatusOK,
		}, nil)

		client := gitlab.NewClient(
			"test_token",
			gitlab.ClientOptionBaseUrl(baseUrl),
			gitlab.ClientOptionHttpClient(httpClient),
		)

		user, err := client.GetUserByID(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, userID, user.ID)
	})

	t.Run("error on getting user", func(t *testing.T) {
		expErr := errors.New("test error")

		httpClient := new(gitlab.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, expErr)

		client := gitlab.NewClient(
			"test_token",
			gitlab.ClientOptionHttpClient(httpClient),
		)

		user, err := client.GetUserByID(context.Background(), 5)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expErr))
		assert.Equal(t, gitlab.User{}, user)
	})
}

func TestClient_GetDiscussion(t *testing.T) {
	t.Run("positive case", func(t *testing.T) {
		var (
			projectID    = 10
			mrID         = 20
			discussionID = "test_discussion"
			baseUrl      = "http://gitlab.test.com/api/v4"
		)

		httpClient := new(gitlab.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
			req, ok := args.Get(0).(*http.Request)
			assert.True(t, ok)

			path := fmt.Sprintf("projects/%d/merge_requests/%d/discussions/%s", projectID, mrID, discussionID)
			assert.Equal(t, baseUrl+"/"+path, req.URL.String())
		}).Return(&http.Response{
			Body: ioutil.NopCloser(
				bytes.NewReader(
					[]byte(fmt.Sprintf("{\"id\": \"%s\"}", discussionID)),
				),
			),
			StatusCode: http.StatusOK,
		}, nil)

		client := gitlab.NewClient(
			"test_token",
			gitlab.ClientOptionBaseUrl(baseUrl),
			gitlab.ClientOptionHttpClient(httpClient),
		)

		discussion, err := client.GetDiscussion(context.Background(), projectID, mrID, discussionID)
		assert.NoError(t, err)
		assert.Equal(t, discussionID, discussion.ID)
	})

	t.Run("error on getting discussion", func(t *testing.T) {
		expErr := errors.New("test error")

		httpClient := new(gitlab.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, expErr)

		client := gitlab.NewClient(
			"test_token",
			gitlab.ClientOptionHttpClient(httpClient),
		)

		discussion, err := client.GetDiscussion(context.Background(), 10, 20, "test_discussion")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expErr))
		assert.Equal(t, gitlab.Discussion{}, discussion)
	})
}

func TestClient_GetParticipants(t *testing.T) {
	t.Run("positive case", func(t *testing.T) {
		var (
			projectID          = 10
			mrID               = 20
			discussionID       = "test_discussion"
			baseUrl            = "http://gitlab.test.com/api/v4"
			expParticipantsIDs = map[int]struct{}{
				15: {},
				25: {},
				35: {},
			}
		)

		// preparing response mock
		expResponse := gitlab.Discussion{Notes: []gitlab.Note{}}
		for id := range expParticipantsIDs {
			expResponse.Notes = append(expResponse.Notes, gitlab.Note{
				Author: gitlab.NoteAuthor{ID: id},
			})
		}

		expResponseBytes, err := json.Marshal(expResponse)
		assert.NoError(t, err)

		httpClient := new(gitlab.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
			req, ok := args.Get(0).(*http.Request)
			assert.True(t, ok)

			path := fmt.Sprintf("projects/%d/merge_requests/%d/discussions/%s", projectID, mrID, discussionID)
			assert.Equal(t, baseUrl+"/"+path, req.URL.String())
		}).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader(expResponseBytes)),
			StatusCode: http.StatusOK,
		}, nil)

		client := gitlab.NewClient(
			"test_token",
			gitlab.ClientOptionBaseUrl(baseUrl),
			gitlab.ClientOptionHttpClient(httpClient),
		)

		participants, err := client.GetParticipants(context.Background(), projectID, mrID, discussionID)
		assert.NoError(t, err)
		assert.Equal(t, len(expParticipantsIDs), len(participants))

		for _, participant := range participants {
			_, has := expParticipantsIDs[participant.ID]
			assert.True(t, has)
		}
	})

	t.Run("error on getting discussion", func(t *testing.T) {
		expErr := errors.New("test error")

		httpClient := new(gitlab.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, expErr)

		client := gitlab.NewClient(
			"test_token",
			gitlab.ClientOptionHttpClient(httpClient),
		)

		participants, err := client.GetParticipants(context.Background(), 10, 20, "test_discussion")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expErr))
		assert.Equal(t, []gitlab.NoteAuthor(nil), participants)
	})
}
