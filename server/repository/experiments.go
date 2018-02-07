package repository

import (
	"database/sql"
	"fmt"

	"github.com/src-d/code-annotation/server/model"
)

// Experiments repository
type Experiments struct {
	db *sql.DB
}

// NewExperiments returns a new Experiments repository
func NewExperiments(db *sql.DB) *Experiments {
	return &Experiments{db: db}
}

// getWithQuery builds an Experiment from the given sql QueryRow. If the
// Experiment does not exist, it returns nil, nil
func (repo *Experiments) getWithQuery(queryRow *sql.Row) (*model.Experiment, error) {
	var exp model.Experiment

	err := queryRow.Scan(&exp.ID, &exp.Name, &exp.Description)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("Error getting experiment from the DB: %v", err)
	default:
		return &exp, nil
	}
}

const selectExperimentsSQL = `SELECT * FROM experiments WHERE id=$1`

// GetByID returns the Experiment with the given ID. If the Experiment does not
// exist, it returns nil, nil
func (repo *Experiments) GetByID(id int) (*model.Experiment, error) {
	return repo.getWithQuery(repo.db.QueryRow(selectExperimentsSQL, id))
}
