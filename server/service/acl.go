package service

import (
	"net/http"

	"github.com/src-d/code-annotation/server/model"
	"github.com/src-d/code-annotation/server/repository"
)

// ACL allows to manage permissions for user role
type ACL struct {
	usersRepo *repository.Users
	role      model.Role
}

// NewACL creates new ACL service
func NewACL(usersRepo *repository.Users, role model.Role) *ACL {
	return &ACL{
		usersRepo: usersRepo,
		role:      role,
	}
}

// Middleware provides middleware which allows requests only for user only with selected role
func (s *ACL) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := GetUserID(r.Context())

		user, err := s.usersRepo.GetByID(userID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user == nil || user.Role != s.role {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
