package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/src-d/code-annotation/server/handler"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"
)

type appConfig struct {
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT" default:"8080"`
	UIDomain string `envconfig:"UI_DOMAIN" default:"http://127.0.0.1:8080"`
}

func main() {
	// main configuration
	var conf appConfig
	envconfig.MustProcess("", &conf)

	// create repos
	userRepo := &repository.Users{}

	// create services
	var oauthConfig service.OAuthConfig
	envconfig.MustProcess("oauth", &oauthConfig)
	oauth := service.NewOAuth(oauthConfig.ClientID, oauthConfig.ClientSecret)

	var jwtConfig service.JWTConfig
	envconfig.MustProcess("jwt", &jwtConfig)
	jwt := service.NewJWT(jwtConfig.SigningKey)

	// cors options
	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Location", "Authorization"},
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.New(corsOptions).Handler)

	r.Get("/login", handler.Login(oauth))
	r.Get("/oauth-callback", handler.OAuthCallback(oauth, jwt, userRepo, conf.UIDomain))

	r.Route("/api", func(r chi.Router) {
		r.Use(jwt.Middleware)
		r.Get("/me", handler.Me(userRepo))
	})

	r.Route("/experiments/{experimentId}", func(r chi.Router) {
		r.Get("/", handler.GetExperimentDetails)

		r.Route("/assignments", func(r chi.Router) {
			r.Get("/", handler.GetAssignmentsForUserExperiment)
			r.Put("/{assignmentId}", handler.SaveAssignment)
		})

		r.Get("/filePairs/{pairId}", handler.GetFilePairDetails)
	})

	r.Get("/*", handler.FrontendStatics)

	//TODO: logging
	logrus.Info("running...")
	http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), r)
}
