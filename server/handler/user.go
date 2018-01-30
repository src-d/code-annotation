package handler

import (
	"fmt"
	"net/http"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/src-d/code-annotation/server/service"
)

// Me handler returns information about current user
func Me(usersRepo *repository.Users) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		uID := service.GetUserID(r.Context())
		if uID == 0 {
			return nil, fmt.Errorf("no user id in context")
		}

		u, err := usersRepo.Get(uID)
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusNotFound, "user not found")
		}

		return serializer.NewUserResponse(u), nil
	}
}
