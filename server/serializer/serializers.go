package serializer

import (
	"net/http"
	"strings"

	"github.com/src-d/code-annotation/server/model"
)

// HTTPError defines an Error message as it will be written in the http.Response
type HTTPError interface {
	error
	StatusCode() int
}

// Response encapsulate the content of an http.Response
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Errors []HTTPError `json:"errors,omitempty"`
}

type httpError struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Details string `json:"details,omitempty"`
}

// StatusCode returns the Status of the httpError
func (e httpError) StatusCode() int {
	return e.Status
}

// Error returns the string content of the httpError
func (e httpError) Error() string {
	if msg := e.Title; msg != "" {
		return msg
	}

	if msg := http.StatusText(e.Status); msg != "" {
		return msg
	}

	return http.StatusText(http.StatusInternalServerError)
}

// NewHTTPError returns an Error
func NewHTTPError(statusCode int, msg ...string) HTTPError {
	return httpError{Status: statusCode, Title: strings.Join(msg, " ")}
}

func newResponse(c interface{}) *Response {
	if c == nil {
		return &Response{
			Status: http.StatusNoContent,
		}
	}

	return &Response{
		Status: http.StatusOK,
		Data:   c,
	}
}

// NewEmptyResponse returns an empty Response
func NewEmptyResponse() *Response {
	return &Response{}
}

type experimentResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewExperimentResponse returns a Response for the passed Experiment
func NewExperimentResponse(e *model.Experiment) *Response {
	return newResponse(experimentResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
	})
}

type assignmentResponse struct {
	ID           int     `json:"id"`
	UserID       int     `json:"userId"`
	PairID       int     `json:"pairId"`
	ExperimentID int     `json:"experimentId"`
	Answer       *string `json:"answer"`
	Duration     int     `json:"duration"`
}

// NewAssignmentsResponse returns a Response for the passed Assignment
func NewAssignmentsResponse(as []*model.Assignment) *Response {
	assignments := make([]assignmentResponse, len(as))
	for i, a := range as {
		var answer *string

		if a.Answer.Valid {
			answer = &a.Answer.String
		}

		assignments[i] = assignmentResponse{a.ID, a.UserID, a.PairID,
			a.ExperimentID, answer, a.Duration}
	}

	return newResponse(assignments)
}

type filePairResponse struct {
	ID   int    `json:"id"`
	Diff string `json:"diff"`
}

// NewFilePairResponse returns a Response for the given FilePair
func NewFilePairResponse(fp *model.FilePair) *Response {
	return newResponse(filePairResponse{fp.ID, fp.Diff})
}

type userResponse struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarURL"`
}

// NewUserResponse returns a Response for the passed User
func NewUserResponse(u *model.User) *Response {
	return newResponse(userResponse{u.ID, u.Login, u.Username, u.AvatarURL})
}

type countResponse struct {
	Count int `json:"count"`
}

// NewCountResponse returns a Response for the total of a count
func NewCountResponse(c int) *Response {
	return newResponse(countResponse{c})
}
