package handler

import (
	"net/http"
	"strconv"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"

	"github.com/go-chi/chi"
)

func GetFilePairDetails() RequestProcessFunc {
	return func(r *http.Request) (resp *serializer.Response, errs []*serializer.Error) {
		pairId, err := strconv.Atoi(chi.URLParam(r, "pairId"))
		if err != nil {
			errs = append(errs, &serializer.Error{500, "Wrong format in file-pair ID sent", err.Error()})
			return
		}

		pairFiles, err := repository.GetFilePairFor(pairId)
		if err != nil {
			errs = append(errs, &serializer.Error{404, "No file-pair found", err.Error()})
			return
		}

		resp = serializer.NewFilePairsResponse(pairFiles)
		return
	}
}
