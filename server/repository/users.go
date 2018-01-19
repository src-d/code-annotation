package repository

import (
	"fmt"

	"github.com/src-d/code-annotation/server/model"
)

type Users struct {
	users []*model.User
}

func (r *Users) Create(m *model.User) error {
	for _, u := range r.users {
		if u.ID == m.ID {
			return nil
		}
	}
	r.users = append(r.users, m)
	return nil
}

func (r *Users) Get(id int) (*model.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user with id %d not found", id)
}
