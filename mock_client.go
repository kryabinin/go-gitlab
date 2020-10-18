// Code generated by mockery v1.0.0. DO NOT EDIT.

package gitlab

import context "context"
import mock "github.com/stretchr/testify/mock"

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

// GetDiscussion provides a mock function with given fields: ctx, projectID, mrID, discussionID
func (_m *MockClient) GetDiscussion(ctx context.Context, projectID int, mrID int, discussionID string) (Discussion, error) {
	ret := _m.Called(ctx, projectID, mrID, discussionID)

	var r0 Discussion
	if rf, ok := ret.Get(0).(func(context.Context, int, int, string) Discussion); ok {
		r0 = rf(ctx, projectID, mrID, discussionID)
	} else {
		r0 = ret.Get(0).(Discussion)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int, string) error); ok {
		r1 = rf(ctx, projectID, mrID, discussionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetParticipants provides a mock function with given fields: ctx, projectID, mrID, discussionID
func (_m *MockClient) GetParticipants(ctx context.Context, projectID int, mrID int, discussionID string) ([]NoteAuthor, error) {
	ret := _m.Called(ctx, projectID, mrID, discussionID)

	var r0 []NoteAuthor
	if rf, ok := ret.Get(0).(func(context.Context, int, int, string) []NoteAuthor); ok {
		r0 = rf(ctx, projectID, mrID, discussionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]NoteAuthor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int, string) error); ok {
		r1 = rf(ctx, projectID, mrID, discussionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByID provides a mock function with given fields: ctx, userID
func (_m *MockClient) GetUserByID(ctx context.Context, userID int) (User, error) {
	ret := _m.Called(ctx, userID)

	var r0 User
	if rf, ok := ret.Get(0).(func(context.Context, int) User); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
