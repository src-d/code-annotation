package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"
)

// Login handler redirects user to oauth provider
func Login(oAuth *service.OAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := oAuth.MakeAuthURL()
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// OAuthCallback makes exchange with oauth provider, gets&creates user and redirects to index page with JWT token
func OAuthCallback(oAuth *service.OAuth, jwt *service.JWT, userRepo *repository.Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := oAuth.ValidateState(r.FormValue("state")); err != nil {
			writeResponse(w, respErr(http.StatusBadRequest, err.Error()))
			return
		}

		code := r.FormValue("code")
		user, err := oAuth.GetUser(r.Context(), code)
		if err != nil {
			logrus.Errorf("oauth get user error: %s", err)
			// FIXME can it be not server error? for wrong code
			writeResponse(w, respInternalErr())
			return
		}

		// FIXME with real repo we need to check does user exists already or not
		if err := userRepo.Create(user); err != nil {
			logrus.Errorf("can't create user: %s", err)
			writeResponse(w, respInternalErr())
			return
		}

		token, err := jwt.MakeToken(user)
		if err != nil {
			logrus.Errorf("make jwt token error: %s", err)
			writeResponse(w, respInternalErr())
			return
		}
		url := "/?token=" + token
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
