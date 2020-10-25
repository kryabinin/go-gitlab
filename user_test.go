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

func TestClient_GetUsersByIDs(t *testing.T) {
	t.Run("positive case", func(t *testing.T) {
		for _, concurrency := range []int{0, 1, 2, 3, 4} {
			var (
				im = map[int]struct{}{
					5:  {},
					10: {},
					15: {},
				}
				ids = make([]int, 0, len(im))
			)

			httpClient := new(gitlab.MockHTTPClient)
			for id := range im {
				ids = append(ids, id)

				httpClient.On("Do", mock.AnythingOfType("*http.Request")).
					Return(&http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(fmt.Sprintf("{\"id\": %d}", id)))),
						StatusCode: http.StatusOK,
					}, nil).
					Once()
			}

			client := gitlab.NewClient(
				"test_token",
				gitlab.WithConcurrency(concurrency),
				gitlab.WithHttpClient(httpClient),
			)

			users, err := client.GetUsersByIDs(context.Background(), ids)
			assert.NoError(t, err)
			assert.Equal(t, len(ids), len(users))
			for _, user := range users {
				_, ok := im[user.ID]
				assert.True(t, ok)
			}
		}
	})

	t.Run("error on getting user", func(t *testing.T) {
		expErr := errors.New("test error")

		for _, concurrency := range []int{0, 1, 2, 3, 4} {
			var (
				im = map[int]error{
					6:  nil,
					11: expErr,
					16: nil,
				}
				ids = make([]int, 0, len(im))
			)

			httpClient := new(gitlab.MockHTTPClient)
			for id, err := range im {
				ids = append(ids, id)

				httpClient.On("Do", mock.AnythingOfType("*http.Request")).
					Return(&http.Response{
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(fmt.Sprintf("{\"id\": %d}", id)))),
						StatusCode: http.StatusOK,
					}, err).
					Once()
			}

			client := gitlab.NewClient(
				"test_token",
				gitlab.WithConcurrency(concurrency),
				gitlab.WithHttpClient(httpClient),
			)

			users, err := client.GetUsersByIDs(context.Background(), ids)
			assert.True(t, errors.Is(err, expErr))
			assert.Equal(t, []gitlab.User(nil), users)
		}
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
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(fmt.Sprintf("{\"id\": %d}", userID)))),
			StatusCode: http.StatusOK,
		}, nil)

		client := gitlab.NewClient(
			"test_token",
			gitlab.WithBaseUrl(baseUrl),
			gitlab.WithHttpClient(httpClient),
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
			gitlab.WithHttpClient(httpClient),
		)

		user, err := client.GetUserByID(context.Background(), 5)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expErr))
		assert.Equal(t, gitlab.User{}, user)
	})

	t.Run("error on unmarshal response", func(t *testing.T) {
		httpClient := new(gitlab.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("{"))),
			StatusCode: http.StatusOK,
		}, nil)

		client := gitlab.NewClient(
			"test_token",
			gitlab.WithHttpClient(httpClient),
		)

		discussion, err := client.GetUserByID(context.Background(), 10)
		assert.Error(t, err)
		assert.Equal(t, gitlab.User{}, discussion)
	})
}
