package handler

import (
	"net/http"

	"github.com/src-d/code-annotation/server/serializer"
)

// GetAssignmentsForUserExperiment returns a function that returns a *serializer.Response
// with the assignments for the logged user and a passed experiment
// if these assignments do not already exist, they are created in advance
func GetAssignmentsForUserExperiment() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		return nil, serializer.NewHTTPError(http.StatusNotImplemented)
	}
}

type assignmentRequest struct {
	Answer   string `json:"answer"`
	Duration int    `json:"duration"`
}

// SaveAssignment returns a function that saves the user answers as passed in the body request
func SaveAssignment() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		return nil, serializer.NewHTTPError(http.StatusNotImplemented)
	}
}
