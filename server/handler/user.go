package handler

import (
	"net/http"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"
)

// Me handler returns information about current user
func Me(usersRepo *repository.Users) http.HandlerFunc {
	return render(func(r *http.Request) response {
		uID := service.GetUserID(r.Context())
		if uID == 0 {
			return respErr(http.StatusInternalServerError, "no user id in context")
		}
		u, err := usersRepo.Get(uID)
		if err != nil {
			return respErr(http.StatusNotFound, "user not found")
		}
		return respOK(u)
	})
}
