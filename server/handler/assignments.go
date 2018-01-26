package handler

import (
	"net/http"

	"github.com/src-d/code-annotation/server/serializer"
)

func GetAssignmentsForUserExperiment() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		return nil, serializer.NewHTTPError(http.StatusNotImplemented)
	}
}

type assignmentRequest struct {
	Answer   string `json:"answer"`
	Duration int    `json:"duration"`
}

func SaveAssignment() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		return nil, serializer.NewHTTPError(http.StatusNotImplemented)
	}
}
