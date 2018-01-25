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
func Login(oAuth *service.OAuth) RequestProcessFunc {
	return func(r *http.Request) (resp *serializer.Response, errs []*serializer.Error) {
		url := oAuth.MakeAuthURL()
		resp = serializer.NewTemporalRedirect(url)
		return
	}
}

// OAuthCallback makes exchange with oauth provider, gets&creates user and redirects to index page with JWT token
func OAuthCallback(
	oAuth *service.OAuth,
	jwt *service.JWT,
	userRepo *repository.Users,
	uiDomain string,
	logger logrus.FieldLogger,
) RequestProcessFunc {
	return func(r *http.Request) (resp *serializer.Response, errs []*serializer.Error) {
		if err := oAuth.ValidateState(r.FormValue("state")); err != nil {
			errs = append(errs, &serializer.Error{500, "state could not be validated", err.Error()})
			return
		}

		code := r.FormValue("code")
		user, err := oAuth.GetUser(r.Context(), code)
		if err != nil {
			// FIXME can it be not server error? for wrong code
			errs = append(errs, &serializer.Error{500, "oauth get user error", err.Error()})
			return
		}

		// FIXME with real repo we need to check does user exists already or not
		if err := userRepo.Create(user); err != nil {
			errs = append(errs, &serializer.Error{500, "can't create user", err.Error()})
			return
		}

		token, err := jwt.MakeToken(user)
		if err != nil {
			errs = append(errs, &serializer.Error{500, "make jwt token error", err.Error()})
			return
		}

		url := fmt.Sprintf("%s/?token=%s", uiDomain, token)
		resp = serializer.NewTemporalRedirect(url)
		return
	}
}
