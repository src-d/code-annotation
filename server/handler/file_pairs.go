package handler

import (
	"net/http"

	"github.com/src-d/code-annotation/server/serializer"
)

func GetFilePairDetails() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		return nil, serializer.NewHTTPError(http.StatusNotImplemented)
	}
}
