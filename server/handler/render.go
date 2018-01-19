package handler

import (
	"encoding/json"
	"net/http"
)

func respOK(d interface{}) response {
	return response{
		statusCode: http.StatusOK,
		Data:       d,
	}
}

func respErr(statusCode int, msg string) response {
	return response{
		statusCode: statusCode,
		Errors:     []errObj{errObj{Title: msg}},
	}
}

func respInternalErr() response {
	return respErr(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

type errObj struct {
	Title string `json:"title"`
}

type response struct {
	statusCode int
	Data       interface{} `json:"data,omitempty"`
	Errors     []errObj    `json:"errors,omitempty"`
}

type renderFunc func(r *http.Request) response

func render(fn renderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, fn(r))
	}
}

func writeResponse(w http.ResponseWriter, resp response) {
	w.WriteHeader(resp.statusCode)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
