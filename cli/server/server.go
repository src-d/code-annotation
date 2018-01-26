package main

import (
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"

	"github.com/src-d/code-annotation/server"
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

	// loger
	logger := service.NewLogger()

	// start the router
	router := server.Router(logger, jwt, oauth, conf.UIDomain, userRepo, "build")
	logger.Info("running...")
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), router)
	logger.Fatal(err)
}
