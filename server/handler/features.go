package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
)

// GetFeatures returns a function that returns a *serializer.Response
// with the list of features for blobId
func GetFeatures(repo *repository.Features) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		blobID := chi.URLParam(r, "blobId")

		// in the future it should take file by blobID from DB
		// and make API request to ML system
		features, err := repo.GetAll(blobID)
		if err != nil {
			return nil, err
		}

		return serializer.NewFeaturesResponse(features), nil
	}
}
