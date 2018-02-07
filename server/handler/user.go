package handler

import (
	"net/http"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/src-d/code-annotation/server/service"
)

// Me handler returns a function that returns a *serializer.Response
// with the information about the current user
func Me(usersRepo *repository.Users) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		userID, err := service.GetUserID(r.Context())
		if err != nil {
			return nil, err
		}

		u, err := usersRepo.GetByID(userID)
		if err != nil {
			return nil, err
		}

		if u == nil {
			return nil, serializer.NewHTTPError(http.StatusNotFound, "user not found")
		}

		return serializer.NewUserResponse(u), nil
	}
}
