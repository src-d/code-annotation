package serializer

import (
	"net/http"

	"github.com/src-d/code-annotation/server/model"
)

//TODO: Responses and Redirects should implement Serializable{Serialize(w)}

type Response struct {
	Status      int         `json:"status"`
	Data        interface{} `json:"data,omitempty"`
	Errors      []*Error    `json:"errors,omitempty"`
	RedirectURL string
}

type Error struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Details string `json:"details,omitempty"`
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

func NewTemporalRedirect(URL string) *Response {
	return &Response{Status: http.StatusTemporaryRedirect, RedirectURL: URL}
}

type ExperimentResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewExperimentResponse(e *model.Experiment) *Response {
	return newResponse(ExperimentResponse{e.ID, e.Name, e.Description})
}

type AssignmentResponse struct {
	ID       int    `json:"id"`
	PairID   int    `json:"pairId"`
	Answer   string `json:"answer"`
	Duration int    `json:"duration"`
}

type AssignmentRequest struct {
	Answer   string `json:"answer"`
	Duration int    `json:"duration"`
}

func NewAssignmentsResponse(as []*model.Assignment) *Response {
	assignments := make([]AssignmentResponse, len(as))
	for i, a := range as {
		assignments[i] = AssignmentResponse{a.ID, a.PairID, string(a.Answer), a.Duration}
	}

	return newResponse(assignments)
}

func NewAssignmentsFromRequest(as *AssignmentRequest) *model.Assignment {
	return &model.Assignment{
		Answer:   as.Answer,
		Duration: as.Duration,
	}
}

type FilePairsResponse struct {
	ID   int    `json:"id"`
	Diff string `json:"diff"`
}

func NewFilePairsResponse(fp *model.FilePairs) *Response {
	return newResponse(FilePairsResponse{fp.ID, fp.Diff})
}

type UserResponse struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarURL"`
}

func NewUserResponse(u *model.User) *Response {
	return newResponse(UserResponse{u.ID, u.Login, u.Username, u.AvatarURL})
}

type CountResponse struct {
	Count int `json:"count"`
}

func NewCountResponse(c int) *Response {
	return newResponse(CountResponse{c})
}
