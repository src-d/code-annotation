package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"

	"github.com/go-chi/chi"
)

// GetFilePairDetails returns a function that returns a *serializer.Response
// with the details of the requested FilePair
func GetFilePairDetails() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		requestedPairID := chi.URLParam(r, "pairId")
		pairID, err := strconv.Atoi(requestedPairID)
		if err != nil {
			return nil, serializer.NewHTTPError(
				http.StatusBadRequest, fmt.Sprintf("wrong format in file-pair ID sent; received %s", requestedPairID),
			)
		}

		pairFiles, err := repository.GetFilePairFor(pairID)
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusNotFound, "no file-pair found")
		}

		return serializer.NewFilePairsResponse(pairFiles), nil
	}
}
