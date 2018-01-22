package handler

import (
	"net/http"
	"os"
)

func FrontendStatics(w http.ResponseWriter, r *http.Request) {
	filepath := "build" + r.URL.Path
	if _, err := os.Stat(filepath); err == nil {
		http.ServeFile(w, r, filepath)
	}
	http.ServeFile(w, r, "build/index.html")
}
