package handler

import (
	"net/http"
	"strconv"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"

	"github.com/go-chi/chi"
	renderer "github.com/go-chi/render"
)

func GetFilePairDetails(w http.ResponseWriter, r *http.Request) {
	pairId, err := strconv.Atoi(chi.URLParam(r, "pairId"))
	if err != nil {
		renderer.Render(w, r, serializer.NewErrors(err))
		return
	}

	pairFiles, err := repository.GetFilePairFor(pairId)
	if err != nil {
		renderer.Render(w, r, serializer.NewErrors(err))
		return
	}

	renderer.Render(w, r, &serializer.ResponseSingle{Data: serializer.NewFilePairs(&pairFiles)})
}
