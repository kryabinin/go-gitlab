// Package gitlab - user
package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// User entity
type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	UserName    string `json:"username"`
	PublicEmail string `json:"public_email"`
}

func getUsersByIDs(parentCtx context.Context, c *client, ids []int) ([]User, error) {
	var wg sync.WaitGroup
	wg.Add(len(ids))

	wgChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgChan)
	}()

	semaphore := make(chan struct{}, c.concurrency)
	defer close(semaphore)

	errChan := make(chan error, cap(semaphore))
	defer close(errChan)

	ctx, cancelFunc := context.WithCancel(parentCtx)
	defer cancelFunc()

	users := make([]User, 0, len(ids))
	for _, id := range ids {
		semaphore <- struct{}{}
		go func(id int) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			if user, err := c.GetUserByID(ctx, id); err != nil {
				errChan <- err
			} else {
				users = append(users, user)
			}
		}(id)
	}

	select {
	case err := <-errChan:
		cancelFunc()
		<-wgChan
		return nil, fmt.Errorf("can't get users from gitlab: %w", err)
	case <-wgChan:
		return users, nil
	}
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
