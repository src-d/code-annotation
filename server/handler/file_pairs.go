package handler

import (
	"net/http"

	"github.com/src-d/code-annotation/server/serializer"
)

// GetFilePairDetails returns a function that returns a *serializer.Response
// with the details of the requested FilePair
func GetFilePairDetails() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		return nil, serializer.NewHTTPError(http.StatusNotImplemented)
	}
}
