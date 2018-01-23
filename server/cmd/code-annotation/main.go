package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
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

	// DB
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}

	// routing
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers.Set("Access-Control-Allow-Origin", "*")
			headers.Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			headers.Set("Access-Control-Allow-Headers", "Location, Authorization,Content-Type")
			if r.Method == "OPTIONS" {
				return
			}
			h.ServeHTTP(w, r)
		})
	})

	r.Get("/login", handler.Login(oauth))
	r.Get("/oauth-callback", handler.OAuthCallback(oauth, jwt, userRepo, conf.UIDomain))

	// protected handlers
	r.Route("/api", func(r chi.Router) {
		r.Use(jwt.Middleware)

		r.Get("/me", handler.Me(userRepo))
		r.Get("/experiments/{expId}/file-pairs/{id}", handler.FilePair(db))
		r.Put("/experiments/{expId}/assignments/{id}", handler.UpdateAssignment(db))
		r.Get("/experiments/{id}/assignments", handler.Assignments(db))
		r.Get("/experiments/{id}", handler.Experiment(db))
	})

	// frontend
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		filepath := "build" + r.URL.Path
		if _, err := os.Stat(filepath); err == nil {
			http.ServeFile(w, r, filepath)
			return
		}
		http.ServeFile(w, r, "build/index.html")
	})

	logrus.Info("running...")
	http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), r)
}
