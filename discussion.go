// Package gitlab - discussion
package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
)

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

func getParticipants(ctx context.Context, c *client, projectID, mrID int, discussionID string) ([]NoteAuthor, error) {
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

func getDiscussion(ctx context.Context, c *client, projectID, mrID int, discussionID string) (Discussion, error) {
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
