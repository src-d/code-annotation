package server

import (
	"database/sql"
	"net/http"

	"github.com/src-d/code-annotation/server/handler"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// Router returns a Handler to serve the code-anotation backend
func Router(
	logger logrus.FieldLogger,
	jwt *service.JWT,
	oauth *service.OAuth,
	uiDomain string,
	db *sql.DB,
	staticsPath string,
) http.Handler {

	// create repos
	userRepo := repository.NewUsers(db)
	experimentRepo := repository.NewExperiments(db)
	assignmentRepo := repository.NewAssignments(db)
	filePairRepo := repository.NewFilePairs(db)

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

	r.Get("/login", handler.Login(oauth))
	r.Get("/oauth-callback", handler.OAuthCallback(oauth, jwt, userRepo, uiDomain, logger))

	r.Route("/api", func(r chi.Router) {
		r.Use(jwt.Middleware)

		r.Get("/me", handler.Get(handler.Me(userRepo)))

		r.Route("/experiments/{experimentId}", func(r chi.Router) {

			r.Get("/", handler.Get(handler.GetExperimentDetails(experimentRepo)))

			r.Route("/assignments", func(r chi.Router) {

				r.Get("/", handler.Get(handler.GetAssignmentsForUserExperiment(assignmentRepo)))
				r.Put("/{assignmentId}", handler.Get(handler.SaveAssignment(assignmentRepo)))
			})

			r.Get("/file-pairs/{pairId}", handler.Get(handler.GetFilePairDetails(filePairRepo)))
		})
	})

	r.Get("/static/*", handler.FrontendStatics(staticsPath, false))
	r.Get("/*", handler.FrontendStatics(staticsPath, true))

	return r
}
