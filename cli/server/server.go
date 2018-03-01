package main

import (
	"fmt"
	"net/http"

	"github.com/src-d/code-annotation/server"
	"github.com/src-d/code-annotation/server/dbutil"
	"github.com/src-d/code-annotation/server/handler"
	"github.com/src-d/code-annotation/server/service"

	"github.com/kelseyhightower/envconfig"
)

// version will be replaced automatically by the CI build.
// See https://github.com/src-d/ci/blob/v1/Makefile.main#L56
var version = "dev"

type appConfig struct {
	Env          string `envconfig:"ENV" default:"production"`
	Host         string `envconfig:"HOST" default:"0.0.0.0"`
	Port         int    `envconfig:"PORT" default:"8080"`
	ServerURL    string `envconfig:"SERVER_URL"`
	DBConn       string `envconfig:"DB_CONNECTION" default:"sqlite:///var/code-annotation/internal.db"`
	ExportsPath  string `envconfig:"EXPORTS_PATH" default:"./exports"`
	GaTrackingID string `envconfig:"GA_TRACKING_ID" required:"false"`
}

func main() {
	// main configuration
	var conf appConfig
	envconfig.MustProcess("CAT", &conf)
	if conf.ServerURL == "" {
		conf.ServerURL = fmt.Sprintf("//%s:%d", conf.Host, conf.Port)
	}

	// loger
	logger := service.NewLogger(conf.Env)

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
	envconfig.MustProcess("CAT_OAUTH", &oauthConfig)
	oauth := service.NewOAuth(
		oauthConfig.ClientID, oauthConfig.ClientSecret,
		oauthConfig.RestrictAccess, oauthConfig.RestrictRequesterAccess,
	)

	var jwtConfig service.JWTConfig
	envconfig.MustProcess("CAT_JWT", &jwtConfig)
	jwt := service.NewJWT(jwtConfig.SigningKey)

	static := handler.NewStatic("build", conf.ServerURL, conf.GaTrackingID)

	// start the router
	router := server.Router(logger, jwt, oauth, static, &db, conf.ExportsPath, version)
	logger.Info("running...")
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), router)
	logger.Fatal(err)
}
