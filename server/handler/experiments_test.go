package handler_test

import (
	"database/sql"
	"net/http"
	"strings"
	"testing"

	"github.com/src-d/code-annotation/server/dbutil"
	"github.com/src-d/code-annotation/server/handler"
	"github.com/src-d/code-annotation/server/model"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/stretchr/testify/assert"
)

func TestCreateExperiment(t *testing.T) {
	assert := assert.New(t)

	db := testDB()
	repo := repository.NewExperiments(db)
	handler := handler.CreateExperiment(repo)

	json := `{"name": "new", "description": "test"}`
	req, _ := http.NewRequest("POST", "/experiments", strings.NewReader(json))
	res, err := handler(req)
	assert.Nil(err)

	assert.Equal(res, serializer.NewExperimentResponse(&model.Experiment{
		ID:          1,
		Name:        "new",
		Description: "test",
	}, 0))
}

func testDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	err = dbutil.Bootstrap(dbutil.DB{
		DB:     db,
		Driver: dbutil.Sqlite,
	})
	if err != nil {
		panic(err)
	}
	return db
}
