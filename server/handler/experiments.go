package handler

import (
	"net/http"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
)

// GetExperimentDetails returns a function that returns a *serializer.Response
// with the details of a requested experiment
func GetExperimentDetails(repo *repository.Experiments) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		experimentID, err := urlParamInt(r, "experimentId")
		if err != nil {
			return nil, err
		}

		experiment, err := repo.GetByID(experimentID)
		if err != nil {
			return nil, err
		}

		if experiment == nil {
			return nil, serializer.NewHTTPError(http.StatusNotFound, "no experiment found")
		}

		return serializer.NewExperimentResponse(experiment), nil
	}
}
