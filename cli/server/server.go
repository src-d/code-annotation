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
	Host        string `envconfig:"HOST"`
	Port        int    `envconfig:"PORT" default:"8080"`
	UIDomain    string `envconfig:"UI_DOMAIN" default:"http://127.0.0.1:8080"`
	DBConn      string `envconfig:"DB_CONNECTION" default:"sqlite://./internal.db"`
	ExportsPath string `envconfig:"EXPORTS_PATH" default:"./exports"`
}

func main() {
	// main configuration
	var conf appConfig
	envconfig.MustProcess("", &conf)

	// loger
	logger := service.NewLogger()

	// database
	db, err := dbutil.Open(conf.DBConn, false)
	if err != nil {
		logger.Fatalf("error opening the database: %s", err)
	}
	defer db.Close()

	if err := dbutil.Bootstrap(db); err != nil {
		logger.Fatalf("error bootstrapping the database: %s", err)
	}

	if err := dbutil.Initialize(db); err != nil {
		logger.Fatalf("error initializing the database: %s", err)
	}

	// create services
	var oauthConfig service.OAuthConfig
	envconfig.MustProcess("oauth", &oauthConfig)
	oauth := service.NewOAuth(
		oauthConfig.ClientID, oauthConfig.ClientSecret,
		oauthConfig.RestrictAccess, oauthConfig.RestrictRequesterAccess,
	)

	var jwtConfig service.JWTConfig
	envconfig.MustProcess("jwt", &jwtConfig)
	jwt := service.NewJWT(jwtConfig.SigningKey)

	// start the router
	router := server.Router(logger, jwt, oauth, conf.UIDomain, &db, "build", conf.ExportsPath)
	logger.Info("running...")
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), router)
	logger.Fatal(err)
}
