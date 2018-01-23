package handler

import (
	"net/http"
	"strconv"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"

	"github.com/go-chi/chi"
)

func GetExperimentDetails() RequestProcessFunc {
	return func(r *http.Request) (resp *serializer.Response, errs []*serializer.Error) {
		experimentID, err := strconv.Atoi(chi.URLParam(r, "experimentId"))
		if err != nil {
			errs = append(errs, &serializer.Error{500, "Wrong format in experiment ID sent", err.Error()})
			return
		}

		experiment, err := repository.GetExperimentById(experimentID)
		if err != nil {
			errs = append(errs, &serializer.Error{404, "No experiment found", err.Error()})
			return
		}

		resp = serializer.NewExperimentResponse(experiment)
		return
	}
}
