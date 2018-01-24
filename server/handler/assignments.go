package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"

	"github.com/go-chi/chi"
)

func GetAssignmentsForUserExperiment() RequestProcessFunc {
	return func(r *http.Request) (resp *serializer.Response, errs []*serializer.Error) {
		experimentID, err := strconv.Atoi(chi.URLParam(r, "experimentId"))
		if err != nil {
			errs = append(errs, &serializer.Error{500, "Wrong format in experiment ID sent", err.Error()})
			return
		}

		userID, err := strconv.Atoi("1")
		if err != nil {
			errs = append(errs, &serializer.Error{500, "Wrong format in user ID sent", err.Error()})
			return
		}

		assignments, err := repository.GetAssignmentsFor(userID, experimentID)
		if err == repository.ErrNoAssignmentsInitialized {
			assignments, err = repository.CreateAssignmentsFor(userID, experimentID)
		}
		if err != nil {
			errs = append(errs, &serializer.Error{500, "No assignments retrieved", err.Error()})
			return
		}

		resp = serializer.NewAssignmentsResponse(assignments)
		return
	}
}

func SaveAssignment() RequestProcessFunc {
	return func(r *http.Request) (resp *serializer.Response, errs []*serializer.Error) {
		assignmentId, err := strconv.Atoi(chi.URLParam(r, "assignmentId"))
		if err != nil {
			errs = append(errs, &serializer.Error{500, "Wrong format in assignment ID sent", err.Error()})
			return
		}

		var assignmentRequest serializer.AssignmentRequest
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err == nil {
			err = json.Unmarshal(body, &assignmentRequest)
		}
		if err != nil {
			errs = append(errs, &serializer.Error{500, "Payload could not be read", err.Error()})
			return
		}

		if err := repository.UpdateAssignment(assignmentId, assignmentRequest.Answer, assignmentRequest.Duration); err != nil {
			errs = append(errs, &serializer.Error{500, "Answer could not be saved", err.Error()})
			return
		}

		resp = serializer.NewCountResponse(1)
		return
	}
}
