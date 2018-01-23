package handler

import (
	"net/http"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/src-d/code-annotation/server/service"
)

// Me handler returns information about current user
func Me(usersRepo *repository.Users) RequestProcessFunc {
	return func(r *http.Request) (resp *serializer.Response, errs []*serializer.Error) {
		uID := service.GetUserID(r.Context())
		if uID == 0 {
			errs = append(errs, &serializer.Error{500, "no user id in context", ""})
		}

		u, err := usersRepo.Get(uID)
		if err != nil {
			errs = append(errs, &serializer.Error{404, "user not found", err.Error()})
		}

		resp = serializer.NewUserResponse(u)
		return
	}
}
