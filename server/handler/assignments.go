package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/src-d/code-annotation/server/service"
)

// GetAssignmentsForUserExperiment returns a function that returns a *serializer.Response
// with the assignments for the logged user and a passed experiment
// if these assignments do not already exist, they are created in advance
func GetAssignmentsForUserExperiment(repo *repository.Assignments) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		experimentID, err := urlParamInt(r, "experimentId")
		if err != nil {
			return nil, err
		}

		userID, err := service.GetUserID(r.Context())
		if err != nil {
			return nil, err
		}

		initialized, err := repo.IsInitialized(userID, experimentID)
		if err != nil {
			return nil, err
		}

		if !initialized {
			if _, err = repo.Initialize(userID, experimentID); err != nil {
				return nil, err
			}
		}

		assignments, err := repo.GetAll(userID, experimentID)
		if err != nil {
			return nil, err
		}

		return serializer.NewAssignmentsResponse(assignments), nil
	}
}

type assignmentRequest struct {
	Answer   string `json:"answer"`
	Duration int    `json:"duration"`
}

// SaveAssignment returns a function that saves the user answers as passed in the body request
func SaveAssignment(repo *repository.Assignments) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		assignmentID, err := urlParamInt(r, "assignmentId")
		if err != nil {
			return nil, err
		}

		assignment, err := repo.GetByID(assignmentID)
		if err != nil {
			return nil, err
		}

		if assignment == nil {
			return nil, serializer.NewHTTPError(http.StatusNotFound, "assignment not found")
		}

		userID, err := service.GetUserID(r.Context())
		if err != nil {
			return nil, err
		}

		if userID != assignment.UserID {
			return nil, serializer.NewHTTPError(http.StatusForbidden,
				"logged in user is not the assignment's owner")
		}

		var assignmentRequest assignmentRequest
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err == nil {
			err = json.Unmarshal(body, &assignmentRequest)
		}

		if err != nil {
			return nil, err
		}

		err = repo.Update(assignmentID, assignmentRequest.Answer, assignmentRequest.Duration)
		if err != nil {
			return nil, err
		}

		return serializer.NewCountResponse(1), nil
	}
}

// GetFilePairAnnotations returns a function that returns a *serializer.Response
// with the Annotation results for the given File Pair and Experiment IDs
func GetFilePairAnnotations(repo *repository.Assignments) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		experimentID, err := urlParamInt(r, "experimentId")
		if err != nil {
			return nil, err
		}

		pairID, err := urlParamInt(r, "pairId")
		if err != nil {
			return nil, err
		}

		assignments, err := repo.GetByExperimentPair(experimentID, pairID)
		if err != nil {
			return nil, err
		}

		var responseData serializer.ExpAnnotationResponse

		for _, a := range assignments {
			switch a.AnswerStr() {
			case "yes":
				responseData.Yes++
			case "maybe":
				responseData.Maybe++
			case "no":
				responseData.No++
			case "skip":
				responseData.Skip++
			case "":
				responseData.Unanswered++
			}
		}

		responseData.Total = len(assignments)

		return serializer.NewExpAnnotationsResponse(responseData), nil
	}
}
