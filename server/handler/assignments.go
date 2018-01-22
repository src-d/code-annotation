package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	renderer "github.com/go-chi/render"

	"github.com/src-d/code-annotation/server/model"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
)

func GetAssignmentsForUserExperiment(w http.ResponseWriter, r *http.Request) {
	experimentID, err := strconv.Atoi(chi.URLParam(r, "experimentId"))
	if err != nil {
		renderer.Render(w, r, serializer.NewErrors(err))
		return
	}

	userID, err := strconv.Atoi("1")
	if err != nil {
		renderer.Render(w, r, serializer.NewErrors(err))
		return
	}

	assignments, err := repository.GetAssignmentsFor(userID, experimentID)
	if err != nil {
		//TODO: Create assignments
		renderer.Render(w, r, serializer.NewErrors(err))
		return
	}

	renderer.Render(w, r, &serializer.ResponseMany{Data: serializer.NewAssignments(assignments)})
}

func SaveAssignment(w http.ResponseWriter, r *http.Request) {
	assignmentId, err := strconv.Atoi(chi.URLParam(r, "assignmentId"))
	if err != nil {
		renderer.Render(w, r, serializer.NewErrors(err))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderer.Render(w, r, serializer.NewErrors(err))
		return
	}

	var assignment serializer.Assignment
	if err := json.Unmarshal(body, &assignment); err != nil {
		renderer.Render(w, r, serializer.NewErrors(err))
		return
	}

	assignmentCandidate := serializer.NewAssignmentUpdateRequest(assignment)
	if err := repository.UpdateAssignment(assignmentId, assignmentCandidate); err != nil {
		renderer.Render(w, r, serializer.NewErrors(err))
		return
	}

	//TODO: use chi logger
	log.Printf("\nstored: %#v", model.Assignment{
		ID:       assignmentId,
		Answer:   assignment.Answer,
		Duration: assignment.Duration},
	)

	renderer.Render(w, r, &serializer.ResponseStored{1})
}
