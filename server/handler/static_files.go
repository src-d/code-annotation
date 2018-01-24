package handler

import (
	"net/http"
	"os"
)

// FrontendStatics handles the static-files requests, serving the files under the given staticsPath
// if useIndexFallback is set to true and the requested file does not exist, the staticsPath/index.html will be served
func FrontendStatics(staticsPath string, useIndexFallback bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filepath := staticsPath + r.URL.Path
		if useIndexFallback {
			if _, err := os.Stat(filepath); err != nil {
				filepath = staticsPath + "/index.html"
			}
		}

		http.ServeFile(w, r, filepath)
	}
}
