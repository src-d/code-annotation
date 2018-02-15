package handler

import (
	"bytes"
	"net/http"

	"github.com/src-d/code-annotation/server/assets"
)

// FrontendStatics handles the static-files requests, serving the files under the given staticsPath
// if useIndexFallback is set to true and the requested file does not exist, the staticsPath/index.html will be served
func FrontendStatics(staticsPath string, useIndexFallback bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filepath := staticsPath + r.URL.Path
		b, err := assets.Asset(filepath)
		// fallback on index, will be handled by FE router
		if err != nil {
			filepath = staticsPath + "/index.html"
			b, err = assets.Asset(filepath)
			if err != nil {
				http.NotFound(w, r)
				return
			}
		}
		info, err := assets.AssetInfo(filepath)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		http.ServeContent(w, r, info.Name(), info.ModTime(), bytes.NewReader(b))
	}
}
