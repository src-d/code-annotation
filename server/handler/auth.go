package handler

import (
	"fmt"
	"net/http"

	"github.com/src-d/code-annotation/server/model"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/src-d/code-annotation/server/service"

	"github.com/sirupsen/logrus"
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
			errorText := "The state passed is incorrect or expired"
			write(
				w, r,
				serializer.NewEmptyResponse(),
				serializer.NewHTTPError(http.StatusBadRequest, errorText),
			)
			return
		}

		code := r.FormValue("code")
		if code == "" {
			errorText := r.FormValue("error_description")
			if errorText == "" {
				errorText = "OAuth provided didn't send code in callback"
			}
			write(
				w, r,
				serializer.NewEmptyResponse(),
				serializer.NewHTTPError(http.StatusBadRequest, errorText),
			)
			return
		}

		ghUser, err := oAuth.GetUser(r.Context(), code)
		if err == service.ErrNoAccess {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		if err != nil {
			logger.Errorf("oauth get user error: %s", err)
			// FIXME can it be not server error? for wrong code
			write(w, r, serializer.NewEmptyResponse(), err)
			return
		}

		user, err := userRepo.Get(ghUser.Login)
		if err != nil {
			logger.Error(err)
			write(w, r, serializer.NewEmptyResponse(), err)
			return
		}

		if user == nil {
			user = &model.User{
				Login:     ghUser.Login,
				Username:  ghUser.Username,
				AvatarURL: ghUser.AvatarURL,
				Role:      ghUser.Role,
			}

			err = userRepo.Create(user)
			if err != nil {
				logger.Errorf("can't create user: %s", err)
				write(w, r, serializer.NewEmptyResponse(), err)
				return
			}
		} else {
			user.Username = ghUser.Username
			user.AvatarURL = ghUser.AvatarURL
			user.Role = ghUser.Role

			if err = userRepo.Update(user); err != nil {
				logger.Errorf("can't update user: %s", err)
				write(w, r, serializer.NewEmptyResponse(), err)
				return
			}
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
