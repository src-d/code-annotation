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
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Progress    float32 `json:"progress"`
}

// NewExperimentResponse returns a Response for the passed Experiment
func NewExperimentResponse(e *model.Experiment, progress float32) *Response {
	return newResponse(experimentResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Progress:    progress,
	})
}

// NewExperimentsResponse returns a Response with a list of Experiments
func NewExperimentsResponse(experiments []*model.Experiment, progresses []float32) *Response {
	result := make([]experimentResponse, len(experiments))
	for i, e := range experiments {
		result[i] = experimentResponse{
			ID:          e.ID,
			Name:        e.Name,
			Description: e.Description,
			Progress:    progresses[i],
		}
	}

	return newResponse(result)
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

// ExpAnnotationResponse stores the data needed by NewExpAnnotationsResponse
type ExpAnnotationResponse struct {
	Yes        int `json:"yes"`
	Maybe      int `json:"maybe"`
	No         int `json:"no"`
	Skip       int `json:"skip"`
	Unanswered int `json:"unanswered"`
	Total      int `json:"total"`
}

// NewExpAnnotationsResponse returns a Response for the Experiment Annotation
// results
func NewExpAnnotationsResponse(data ExpAnnotationResponse) *Response {
	return newResponse(data)
}

type filePairResponse struct {
	ID          int     `json:"id"`
	Diff        string  `json:"diff"`
	Score       float64 `json:"score"`
	LeftBlobID  string  `json:"leftBlobId"`
	RightBlobID string  `json:"rightBlobId"`
	LeftLOC     int     `json:"leftLoc"`
	RightLOC    int     `json:"rightLoc"`
}

// NewFilePairResponse returns a Response for the given FilePair
func NewFilePairResponse(fp *model.FilePair, leftLOC, rightLOC int) *Response {
	return newResponse(filePairResponse{
		fp.ID, fp.Diff, fp.Score, fp.Left.BlobID, fp.Right.BlobID, leftLOC, rightLOC})
}

type listFilePairResponse struct {
	ID        int    `json:"id"`
	LeftPath  string `json:"leftPath"`
	RightPath string `json:"rightPath"`
}

// NewListFilePairsResponse returns a Response for the given FilePairs
func NewListFilePairsResponse(fps []*model.FilePair) *Response {
	result := make([]listFilePairResponse, len(fps))
	for i, fp := range fps {
		result[i] = listFilePairResponse{fp.ID, fp.Left.Path, fp.Right.Path}
	}

	return newResponse(result)
}

type userResponse struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarURL"`
	Role      string `json:"role"`
}

// NewUserResponse returns a Response for the passed User
func NewUserResponse(u *model.User) *Response {
	return newResponse(
		userResponse{u.ID, u.Login, u.Username, u.AvatarURL, u.Role.String()})
}

type featureResponse struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

// NewFeaturesResponse returns a Response for the passed Features
func NewFeaturesResponse(fs []*model.Feature) *Response {
	features := make([]featureResponse, len(fs))
	for i, f := range fs {
		features[i] = featureResponse(*f)
	}
	return newResponse(features)
}

type countResponse struct {
	Count int `json:"count"`
}

// NewCountResponse returns a Response for the total of a count
func NewCountResponse(c int) *Response {
	return newResponse(countResponse{c})
}

type versionResponse struct {
	Version string `json:"version"`
}

// NewVersionResponse returns a Response with current version of the server
func NewVersionResponse(version string) *Response {
	return newResponse(versionResponse{version})
}

type tokenResponse struct {
	Token string `json:"token"`
}

// NewTokenResponse returns a Response with a token
func NewTokenResponse(token string) *Response {
	return newResponse(tokenResponse{token})
}
