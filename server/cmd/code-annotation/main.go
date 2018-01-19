package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/src-d/code-annotation/server/handler"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"
)

func main() {
	// create repos
	userRepo := &repository.Users{}

	// create services
	var oauthConfig service.OAuthConfig
	envconfig.MustProcess("oauth", &oauthConfig)
	oauth := service.NewOAuth(oauthConfig.ClientID, oauthConfig.ClientSecret)

	var jwtConfig service.JWTConfig
	envconfig.MustProcess("jwt", &jwtConfig)
	jwt := service.NewJWT(jwtConfig.SigningKey)

	// routing
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers.Set("Access-Control-Allow-Origin", "*")
			headers.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			headers.Set("Access-Control-Allow-Headers", "Location, Authorization")
			if r.Method == "OPTIONS" {
				return
			}
			h.ServeHTTP(w, r)
		})
	})

	r.Get("/login", handler.Login(oauth))
	r.Get("/oauth-callback", handler.OAuthCallback(oauth, jwt, userRepo))

	// protected handlers
	r.Route("/api", func(r chi.Router) {
		r.Use(jwt.Middleware)

		r.Get("/me", handler.Me(userRepo))
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
	http.ListenAndServe(":8080", r)
}
