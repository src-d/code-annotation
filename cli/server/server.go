package main

import (
	"fmt"
	"net/http"

	"github.com/src-d/code-annotation/server"
	"github.com/src-d/code-annotation/server/dbutil"
	"github.com/src-d/code-annotation/server/service"

	"github.com/kelseyhightower/envconfig"
)

type appConfig struct {
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT" default:"8080"`
	UIDomain string `envconfig:"UI_DOMAIN" default:"http://127.0.0.1:8080"`
	DBConn   string `envconfig:"DB_CONNECTION" default:"sqlite://./internal.db"`
}

func main() {
	// main configuration
	var conf appConfig
	envconfig.MustProcess("", &conf)

	// loger
	logger := service.NewLogger()

	// database
	db, err := dbutil.Open(conf.DBConn, true)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// create services
	var oauthConfig service.OAuthConfig
	envconfig.MustProcess("oauth", &oauthConfig)
	oauth := service.NewOAuth(oauthConfig.ClientID, oauthConfig.ClientSecret)

	var jwtConfig service.JWTConfig
	envconfig.MustProcess("jwt", &jwtConfig)
	jwt := service.NewJWT(jwtConfig.SigningKey)

	// start the router
	router := server.Router(logger, jwt, oauth, conf.UIDomain, db.SQLDB(), "build")
	logger.Info("running...")
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), router)
	logger.Fatal(err)
}
