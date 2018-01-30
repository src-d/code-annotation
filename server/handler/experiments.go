package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"

	"github.com/go-chi/chi"
)

// GetExperimentDetails returns a function that returns a *serializer.Response
// with the details of a requested experiment
func GetExperimentDetails() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		requestedExperimentID := chi.URLParam(r, "experimentId")
		experimentID, err := strconv.Atoi(requestedExperimentID)
		if err != nil {
			return nil, serializer.NewHTTPError(
				http.StatusBadRequest, fmt.Sprintf("wrong format in experiment ID sent; received %s", requestedExperimentID),
			)
		}

		experiment, err := repository.GetExperimentByID(experimentID)
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusNotFound, "no experiment found")
		}

		return serializer.NewExperimentResponse(experiment), nil
	}
}
