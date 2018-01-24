package handler

import (
	"net/http"
	"os"
)

// FrontendStatics handles the static-files requests, serving the files under the given path
func FrontendStatics(staticsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filepath := staticsPath + r.URL.Path
		if _, err := os.Stat(filepath); err == nil {
			http.ServeFile(w, r, filepath)
			return
		}

		http.ServeFile(w, r, staticsPath+"/index.html")
	}
}
