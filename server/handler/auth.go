package handler

import (
	"fmt"
	"net/http"

	"github.com/src-d/code-annotation/server/model"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/src-d/code-annotation/server/service"

	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
)

// Login handler redirects user to oauth provider
func Login(oAuth *service.OAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := oAuth.MakeAuthURL(w, r)
		if err != nil {
			lg.RequestLog(r).Warn(err.Error())
			http.Error(w,
				http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError,
			)
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// OAuthCallback makes exchange with oauth provider, gets&creates user and redirects to index page with JWT token
func OAuthCallback(
	oAuth *service.OAuth,
	jwt *service.JWT,
	userRepo *repository.Users,
	logger logrus.FieldLogger,
) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		state := r.URL.Query().Get("state")
		if err := oAuth.ValidateState(r, state); err != nil {
			logger.Warn(err)
			return nil, serializer.NewHTTPError(
				http.StatusPreconditionFailed,
				"The state passed by github is incorrect or expired",
			)
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errorText := r.URL.Query().Get("error_description")
			if errorText == "" {
				errorText = "OAuth provided didn't send code in callback"
			}

			return nil, serializer.NewHTTPError(http.StatusBadRequest, errorText)
		}

		ghUser, err := oAuth.GetUser(r.Context(), code)
		if err == service.ErrNoAccess {
			return nil, serializer.NewHTTPError(
				http.StatusForbidden,
				http.StatusText(http.StatusForbidden),
			)
		}

		if err != nil {
			// FIXME can it be not server error? for wrong code
			return nil, fmt.Errorf("oauth get user error: %s", err)
		}

		user, err := userRepo.Get(ghUser.Login)
		if err != nil {
			return nil, fmt.Errorf("get user from db: %s", err)
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
				return nil, fmt.Errorf("can't create user: %s", err)
			}
		} else {
			user.Username = ghUser.Username
			user.AvatarURL = ghUser.AvatarURL
			user.Role = ghUser.Role

			if err = userRepo.Update(user); err != nil {
				return nil, fmt.Errorf("can't update user: %s", err)
			}
		}

		token, err := jwt.MakeToken(user)
		if err != nil {
			return nil, fmt.Errorf("make jwt token error: %s", err)
		}

		return serializer.NewTokenResponse(token), nil
	}
}
