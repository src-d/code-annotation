package handler

import (
	"fmt"
	"net/http"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/src-d/code-annotation/server/service"
)

// GetExperimentDetails returns a function that returns a *serializer.Response
// with the details of a requested experiment
func GetExperimentDetails(repo *repository.Experiments, assignmentsRepo *repository.Assignments) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		userID, err := service.GetUserID(r.Context())
		if err != nil {
			return nil, err
		}

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

		progress, err := experimentProgress(assignmentsRepo, experiment.ID, userID)
		if err != nil {
			return nil, err
		}

		return serializer.NewExperimentResponse(experiment, progress), nil
	}
}

// GetExperiments returns a function that returns a *serializer.Response
// with the list of existing experiments
func GetExperiments(repo *repository.Experiments, assignmentsRepo *repository.Assignments) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		userID, err := service.GetUserID(r.Context())
		if err != nil {
			return nil, err
		}

		experiments, err := repo.GetAll()
		if err != nil {
			return nil, err
		}

		var progresses []float32
		for _, e := range experiments {
			progress, err := experimentProgress(assignmentsRepo, e.ID, userID)
			if err != nil {
				return nil, err
			}
			progresses = append(progresses, progress)
		}

		return serializer.NewExperimentsResponse(experiments, progresses), nil
	}
}

func experimentProgress(repo *repository.Assignments, experimentID int, userID int) (float32, error) {
	countAll, err := repo.CountUserAssigments(experimentID, userID)
	if err != nil {
		return 0, fmt.Errorf("Error count of assigments from the DB: %v", err)
	}

	countComplete, err := repo.CountCompleteUserAssigments(experimentID, userID)
	if err != nil {
		return 0, fmt.Errorf("Error count of complete assigments from the DB: %v", err)
	}

	if countAll == 0 {
		return 0, nil
	}

	return 100.0 * float32(countComplete) / float32(countAll), nil
}
