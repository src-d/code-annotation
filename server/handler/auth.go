package handler

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/src-d/code-annotation/server/service"
)

// Login handler redirects user to oauth provider
func Login(oAuth *service.OAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := oAuth.MakeAuthURL(w, r)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// OAuthCallback makes exchange with oauth provider, gets&creates user and redirects to index page with JWT token
func OAuthCallback(
	oAuth *service.OAuth,
	jwt *service.JWT,
	userRepo *repository.Users,
	uiDomain string,
	logger logrus.FieldLogger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := oAuth.ValidateState(r, r.FormValue("state")); err != nil {
			write(w, r, serializer.NewEmptyResponse(), serializer.NewHTTPError(http.StatusBadRequest))
			return
		}

		code := r.FormValue("code")
		user, err := oAuth.GetUser(r.Context(), code)
		if err != nil {
			logger.Errorf("oauth get user error: %s", err)
			// FIXME can it be not server error? for wrong code
			write(w, r, serializer.NewEmptyResponse(), err)
			return
		}

		// FIXME with real repo we need to check does user exists already or not
		if err := userRepo.Create(user); err != nil {
			logger.Errorf("can't create user: %s", err)
			write(w, r, serializer.NewEmptyResponse(), err)
			return
		}

		token, err := jwt.MakeToken(user)
		if err != nil {
			logger.Errorf("make jwt token error: %s", err)
			write(w, r, serializer.NewEmptyResponse(), err)
			return
		}
		url := fmt.Sprintf("%s/?token=%s", uiDomain, token)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
