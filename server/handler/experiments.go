package handler

import (
	"net/http"

	"github.com/src-d/code-annotation/server/serializer"
)

// GetExperimentDetails returns a function that returns a *serializer.Response
// with the details of a requested experiment
func GetExperimentDetails() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		return nil, serializer.NewHTTPError(http.StatusNotImplemented)
	}
}
