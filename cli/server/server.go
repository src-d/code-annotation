package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/src-d/code-annotation/server/handler"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/kelseyhightower/envconfig"
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

	// loger
	logger := service.NewLogger()

	// start the router
	router := router(logger, jwt, oauth, conf.UIDomain, userRepo, "build")
	logger.Info("running...")
	http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), router)
}

func GetHandler(rp handler.RequestProcessFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, errors := rp(r)
		handler.Write(w, r, response, errors)
	}
}

func router(
	logger logrus.FieldLogger,
	jwt *service.JWT,
	oauth *service.OAuth,
	redirectUrl string,
	userRepo *repository.Users,
	staticsPath string,
) http.Handler {

	// cors options
	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Location", "Authorization", "Content-Type"},
	}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(cors.New(corsOptions).Handler)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger}))

	r.Get("/login", GetHandler(handler.Login(oauth)))
	r.Get("/oauth-callback", GetHandler(handler.OAuthCallback(oauth, jwt, userRepo, redirectUrl, logger)))

	r.Route("/api", func(r chi.Router) {
		r.Use(jwt.Middleware)

		r.Get("/me", GetHandler(handler.Me(userRepo)))

		r.Route("/experiments/{experimentId}", func(r chi.Router) {

			r.Get("/", GetHandler(handler.GetExperimentDetails()))

			r.Route("/assignments", func(r chi.Router) {

				r.Get("/", GetHandler(handler.GetAssignmentsForUserExperiment()))
				r.Put("/{assignmentId}", GetHandler(handler.SaveAssignment()))
			})

			r.Get("/file-pairs/{pairId}", GetHandler(handler.GetFilePairDetails()))
		})
	})

	r.Get("/*", handler.FrontendStatics(staticsPath))

	return r
}
