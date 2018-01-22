package serializer

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/src-d/code-annotation/server/model"
)

type ResponseErrors struct {
	Errors []Error `json:"errors"`
	status int
}

func (e *ResponseErrors) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.status)
	return nil
}

type ResponseSingle struct {
	Data Identifiable `json:"data"`
}

func (e *ResponseSingle) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ResponseStored struct {
	Data int `json:"data"`
}

func (e *ResponseStored) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ResponseMany struct {
	Data []Identifiable `json:"data"`
}

func (e *ResponseMany) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Identifiable interface {
	id() int
}

type ID struct {
	ID int `json:"id"`
}

func (i *ID) id() int {
	return i.ID
}

type Experiment struct {
	ID
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewExperiment(e *model.Experiment) Identifiable {
	if e == nil {
		return nil
	}

	return &Experiment{ID{e.ID}, e.Name, e.Description}
}

type Assignment struct {
	ID
	PairID   int    `json:"pairId"`
	Answer   string `json:"answer"`
	Duration int    `json:"duration"`
}

func NewAssignments(as []model.Assignment) []Identifiable {
	assignments := make([]Identifiable, len(as))
	for i, a := range as {
		assignments[i] = &Assignment{ID{a.ID}, a.PairID, string(a.Answer), a.Duration}
	}

	return assignments
}

func NewAssignmentUpdateRequest(as Assignment) model.Assignment {
	return model.Assignment{
		Answer:   as.Answer,
		Duration: as.Duration,
	}
}

type FilePairs struct {
	ID
	Diff string `json:"diff"`
}

func NewFilePairs(fp *model.FilePairs) Identifiable {
	if fp == nil {
		return nil
	}

	return &FilePairs{ID{fp.ID}, fp.Diff}
}

type Error struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Details string `json:"details"` //optional
}

func NewErrors(errs ...error) *ResponseErrors {
	errors := make([]Error, len(errs))
	for i, e := range errs {
		errors[i] = Error{500, "unexpected error", e.Error()}
	}

	return &ResponseErrors{errors, 500}
}
