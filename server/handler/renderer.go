package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/src-d/code-annotation/server/serializer"
	"github.com/src-d/code-annotation/server/service"
)

func Get(rp RequestProcessFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, errors := rp(r)
		Write(w, r, response, errors)
	}
}

func Write(w http.ResponseWriter, r *http.Request, response *serializer.Response, errs []*serializer.Error) {
	//TODO: Responses and Redirects should implement Serializable{Serialize(w)}
	if response.RedirectURL != "" {
		http.Redirect(w, r, response.RedirectURL, response.Status)
		return
	}

	status := http.StatusOK

	if len(errs) > 0 {
		status = http.StatusInternalServerError
		response = &serializer.Response{Status: status, Errors: errs}
		for _, e := range errs {
			//TODO: Not all errors should be logged as .Error()
			service.NewLogger().Error(e.Title)
		}
	}

	content, err := json.Marshal(response)
	if err != nil {
		marshallerErrString := "%s; response could not be marshalled; %s"
		err = fmt.Errorf(marshallerErrString, http.StatusText(http.StatusInternalServerError), err.Error())
		service.NewLogger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(content)
}

type RequestProcessFunc func(*http.Request) (*serializer.Response, []*serializer.Error)
