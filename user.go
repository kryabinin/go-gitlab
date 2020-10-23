// Package gitlab - user
package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
)

// User entity
type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	UserName    string `json:"username"`
	PublicEmail string `json:"public_email"`
}

func getUserByID(ctx context.Context, c *client, userID int) (User, error) {
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
