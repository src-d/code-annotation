package repository

import (
	"fmt"

	"github.com/src-d/code-annotation/server/model"
)

// Users is the User repository
type Users struct {
	users []*model.User
}

// Create stores a user into the DB
func (r *Users) Create(m *model.User) error {
	for _, u := range r.users {
		if u.ID == m.ID {
			return nil
		}
	}
	r.users = append(r.users, m)
	return nil
}

// Get returns the User identified by the passed ID.
// If the user does not exist, if returns no user nor error
func (r *Users) Get(id int) (*model.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user with id %d not found", id)
}
