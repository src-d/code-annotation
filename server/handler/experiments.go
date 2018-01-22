package handler

import (
	"net/http"
	"strconv"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"

	"github.com/go-chi/chi"
	renderer "github.com/go-chi/render"
)

func GetExperimentDetails(w http.ResponseWriter, r *http.Request) {
	experimentID, err := strconv.Atoi(chi.URLParam(r, "experimentId"))
	if err != nil {
		renderer.Render(w, r, serializer.NewErrors(err))
		return
	}

	experiment, err := repository.GetExperimentById(experimentID)
	if err != nil {
		renderer.Render(w, r, serializer.NewErrors(err))
		return
	}

	renderer.Render(w, r, &serializer.ResponseSingle{Data: serializer.NewExperiment(&experiment)})
}
