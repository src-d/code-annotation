package server

import (
	"net/http"

	"github.com/src-d/code-annotation/server/dbutil"
	"github.com/src-d/code-annotation/server/handler"
	"github.com/src-d/code-annotation/server/model"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pressly/lg"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// Router returns a Handler to serve the code-anotation backend
func Router(
	logger *logrus.Logger,
	jwt *service.JWT,
	oauth *service.OAuth,
	static *handler.Static,
	dbWrapper *dbutil.DB,
	exportsPath string,
	version string,
) http.Handler {

	db := dbWrapper.SQLDB()

	// create repos
	userRepo := repository.NewUsers(db)
	experimentRepo := repository.NewExperiments(db)
	assignmentRepo := repository.NewAssignments(db)
	filePairRepo := repository.NewFilePairs(db)
	featureRepo := repository.NewFeatures(db)

	// cors options
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Location", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}

	requesterACL := service.NewACL(userRepo, model.Requester)
	export := handler.NewExport(dbWrapper, exportsPath)

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(cors.New(corsOptions).Handler)
	r.Use(lg.RequestLogger(logger))

	r.Get("/login", handler.Login(oauth))
	r.Get("/api/auth", handler.Get(handler.OAuthCallback(oauth, jwt, userRepo, logger)))

	r.Route("/api", func(r chi.Router) {
		r.Use(jwt.Middleware)

		r.Get("/me", handler.Get(handler.Me(userRepo)))

		r.Get("/experiments", handler.Get(handler.GetExperiments(experimentRepo, assignmentRepo)))

		r.Route("/experiments/{experimentId}", func(r chi.Router) {

			r.Get("/", handler.Get(handler.GetExperimentDetails(experimentRepo, assignmentRepo)))

			r.Route("/assignments", func(r chi.Router) {

				r.Get("/", handler.Get(handler.GetAssignmentsForUserExperiment(assignmentRepo)))
				r.Put("/{assignmentId}", handler.Get(handler.SaveAssignment(assignmentRepo)))
			})

			r.Route("/file-pairs", func(r chi.Router) {
				r.Use(requesterACL.Middleware)

				r.Get("/", handler.Get(handler.GetFilePairs(filePairRepo)))
				r.Get("/{pairId}/annotations", handler.Get(handler.GetFilePairAnnotations(assignmentRepo)))
			})

			r.Get("/file-pairs/{pairId}", handler.Get(handler.GetFilePairDetails(filePairRepo)))
		})

		r.Route("/features", func(r chi.Router) {
			r.Use(requesterACL.Middleware)

			r.Get("/{blobId}", handler.Get(handler.GetFeatures(featureRepo)))
		})

		r.Route("/exports", func(r chi.Router) {
			r.Use(requesterACL.Middleware)

			r.Get("/", handler.Get(export.List))
			r.Post("/", handler.Get(export.Create))
			r.Get("/{filename}/download", export.Download)
		})
	})

	r.Get("/version", handler.Get(handler.Version(version)))

	r.Get("/static/*", static.ServeHTTP)
	r.Get("/*", static.ServeHTTP)

	return r
}
