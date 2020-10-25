// Package gitlab - client
package gitlab

//go:generate mockery -case=underscore -inpkg -name=Client
//go:generate mockery -case=underscore -inpkg -name=HTTPClient

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
)

const defaultBaseUrl = "http://gitlab.com/api/v4"

var defaultConcurrency = runtime.NumCPU()

// HTTPClient interface to replace default http client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type (
	// Client provides api to work with gitlab entities
	Client interface {
		// GetUsersByIDs returns list of users by ids
		GetUsersByIDs(ctx context.Context, ids []int) ([]User, error)

		// GetUserByID returns single user by id
		GetUserByID(ctx context.Context, userID int) (User, error)

		// GetDiscussion returns discussion data by project id, merge request id and discussion id
		GetDiscussion(ctx context.Context, projectID, mrID int, discussionID string) (Discussion, error)

		// GetParticipants returns all participants from discussion (by project id, merge request id and discussion id)
		GetParticipants(ctx context.Context, projectID, mrID int, discussionID string) ([]NoteAuthor, error)

		// SendRequest send http request to gitlab
		SendRequest(ctx context.Context, method string, path string, data []byte) ([]byte, error)
	}

	client struct {
		token       string
		baseUrl     string
		concurrency int
		httpClient  HTTPClient
	}
)

// NewClient is client constructor
func NewClient(token string, opts ...ClientOption) Client {
	c := &client{
		token:       token,
		baseUrl:     defaultBaseUrl,
		concurrency: defaultConcurrency,
		httpClient:  &http.Client{},
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	return c
}

// GetParticipants implementation
func (c *client) GetParticipants(ctx context.Context, projectID, mrID int, discussionID string) ([]NoteAuthor, error) {
	return getParticipants(ctx, c, projectID, mrID, discussionID)
}

// GetDiscussion implementation
func (c *client) GetDiscussion(ctx context.Context, projectID, mrID int, discussionID string) (Discussion, error) {
	return getDiscussion(ctx, c, projectID, mrID, discussionID)
}

// GetUsersByIDs implementation
func (c *client) GetUsersByIDs(ctx context.Context, ids []int) ([]User, error) {
	return getUsersByIDs(ctx, c, ids)
}

// GetUserByID implementation
func (c *client) GetUserByID(ctx context.Context, id int) (User, error) {
	return getUserByID(ctx, c, id)
}

func (c *client) get(ctx context.Context, path string) ([]byte, error) {
	return c.SendRequest(ctx, http.MethodGet, path, nil)
}

func (c *client) SendRequest(ctx context.Context, method string, path string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(method, c.baseUrl+"/"+path, bytes.NewReader(data))
	if nil != err {
		return nil, fmt.Errorf("can't create http request: %w", err)
	}

	req = req.WithContext(ctx)

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Private-Token", c.token)

	var resp *http.Response
	if resp, err = c.httpClient.Do(req); nil != err {
		return nil, fmt.Errorf("can't send http request: %w", err)
	}

	if http.StatusOK != resp.StatusCode {
		return nil, fmt.Errorf("gitlab respond with %d status code", resp.StatusCode)
	}

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); nil != err {
		return nil, fmt.Errorf("can't read response body: %w", err)
	}

	return body, nil
}
