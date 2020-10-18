// Package gitlab - client
package gitlab

//go:generate mockery -case=underscore -inpkg -name=Client
//go:generate mockery -case=underscore -inpkg -name=HTTPClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const defaultBaseUrl = "http://gitlab.com/api/v4"

// HTTPClient interface to replace default http client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type (
	// Client provides api to work with gitlab entities
	Client interface {
		// GetUserByID returns user data by id
		GetUserByID(ctx context.Context, userID int) (User, error)

		// GetDiscussion returns discussion data by project id, merge request id and discussion id
		GetDiscussion(ctx context.Context, projectID, mrID int, discussionID string) (Discussion, error)

		// GetParticipants returns all participants from discussion (by project id, merge request id and discussion id)
		GetParticipants(ctx context.Context, projectID, mrID int, discussionID string) ([]NoteAuthor, error)

		// SendRequest send http request to gitlab
		SendRequest(ctx context.Context, method string, path string, data []byte) ([]byte, error)
	}

	client struct {
		token      string
		baseUrl    string
		httpClient HTTPClient
	}
)

// User entity
type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	UserName    string `json:"username"`
	PublicEmail string `json:"public_email"`
}

type (
	// Discussion entity
	Discussion struct {
		ID             string `json:"id"`
		IndividualNote bool   `json:"individual_note"`
		Notes          []Note `json:"notes"`
	}

	// Note (comment) entity
	Note struct {
		ID           int        `json:"id"`
		Type         string     `json:"type"`
		Body         string     `json:"body"`
		Author       NoteAuthor `json:"author"`
		CreatedAt    string     `json:"created_at"`
		UpdatedAt    string     `json:"updated_at"`
		System       bool       `json:"system"`
		NoteableID   int        `json:"noteable_id"`
		NoteableType string     `json:"noteable_type"`
		Position     Position   `json:"position"`
		Resolvable   bool       `json:"resolvable"`
		Resolved     bool       `json:"resolved"`
		NoteableIID  int        `json:"noteable_iid"`
	}

	// NoteAuthor entity
	NoteAuthor struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		UserName  string `json:"username"`
		State     string `json:"state"`
		AvatarUrl string `json:"avatar_url"`
		WebUrl    string `json:"web_url"`
	}

	// Position entity
	Position struct {
		BaseSha      string `json:"base_sha"`
		StartSha     string `json:"start_sha"`
		HeadSha      string `json:"head_sha"`
		OldPath      string `json:"old_path"`
		NewPath      string `json:"new_path"`
		PositionType string `json:"position_type"`
		OldLine      int    `json:"old_line"`
		NewLine      int    `json:"new_line"`
	}
)

// NewClient is client constructor
func NewClient(token string, opts ...ClientOption) Client {
	c := &client{
		token:      token,
		baseUrl:    defaultBaseUrl,
		httpClient: &http.Client{},
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	return c
}

// GetParticipants implementation
func (c *client) GetParticipants(ctx context.Context, projectID, mrID int, discussionID string) ([]NoteAuthor, error) {
	discussion, err := c.GetDiscussion(ctx, projectID, mrID, discussionID)
	if err != nil {
		return nil, fmt.Errorf("can't get discussion from gitlab: %w", err)
	}

	saved := map[int]struct{}{}
	participants := make([]NoteAuthor, 0)

	for _, note := range discussion.Notes {
		if _, has := saved[note.Author.ID]; !has {
			participants = append(participants, note.Author)
			saved[note.Author.ID] = struct{}{}
		}
	}

	return participants, nil
}

// GetDiscussion implementation
func (c *client) GetDiscussion(ctx context.Context, projectID, mrID int, discussionID string) (Discussion, error) {
	url := fmt.Sprintf("projects/%d/merge_requests/%d/discussions/%s", projectID, mrID, discussionID)
	resp, err := c.get(ctx, url)
	if err != nil {
		return Discussion{}, err
	}

	var discussion Discussion
	if err = json.Unmarshal(resp, &discussion); err != nil {
		return Discussion{}, fmt.Errorf("can't unmarshal discussion data: %w", err)
	}

	return discussion, nil
}

// GetUserByID implementation
func (c *client) GetUserByID(ctx context.Context, userID int) (User, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/users/%d", userID))
	if err != nil {
		return User{}, err
	}

	var user User
	if err = json.Unmarshal(resp, &user); err != nil {
		return User{}, fmt.Errorf("can't unmarshal user data: %w", err)
	}

	return user, nil
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
