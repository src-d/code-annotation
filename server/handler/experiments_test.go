package handler_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/src-d/code-annotation/server/handler"
	"github.com/src-d/code-annotation/server/model"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/stretchr/testify/assert"
)

func TestCreateExperiment(t *testing.T) {
	assert := assert.New(t)

	db := testDB()
	repo := repository.NewExperiments(db.DB)
	handler := handler.CreateExperiment(repo)

	json := `{"name": "new", "description": "test"}`
	req, _ := http.NewRequest("POST", "/experiments", strings.NewReader(json))
	res, err := handler(req)
	assert.Nil(err)

	assert.Equal(serializer.NewExperimentResponse(&model.Experiment{
		ID:          2,
		Name:        "new",
		Description: "test",
	}, 0), res)
}

func TestUpdateExperiment(t *testing.T) {
	assert := assert.New(t)

	db := testDB()
	repo := repository.NewExperiments(db.DB)
	assignmentsRepo := repository.NewAssignments(db.DB)
	handler := handler.UpdateExperiment(repo, assignmentsRepo)

	json := `{"name": "new", "description": "test"}`
	req, _ := http.NewRequest("PUT", "/experiments/1", strings.NewReader(json))
	req = chiRequest(req, map[string]string{"experimentId": "1"})
	req = reqWithUser(req, 1)
	res, err := handler(req)
	assert.Nil(err)

	assert.Equal(serializer.NewExperimentResponse(&model.Experiment{
		ID:          1,
		Name:        "new",
		Description: "test",
	}, 0), res)

	req, _ = http.NewRequest("PUT", "/experiments/2", strings.NewReader(json))
	req = chiRequest(req, map[string]string{"experimentId": "2"})
	req = reqWithUser(req, 1)
	res, err = handler(req)
	assert.Nil(res)
	assert.Equal(serializer.NewHTTPError(http.StatusNotFound, "no experiment found"), err)
}
