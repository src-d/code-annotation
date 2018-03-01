package repository

import (
	"database/sql"
	"fmt"

	"github.com/src-d/code-annotation/server/model"
)

// Experiments repository
type Experiments struct {
	db          *sql.DB
	assignments *Assignments
}

// NewExperiments returns a new Experiments repository
func NewExperiments(db *sql.DB) *Experiments {
	return &Experiments{db: db, assignments: NewAssignments(db)}
}

// getWithQuery builds an Experiment from the given sql QueryRow. If the
// Experiment does not exist, it returns nil, nil
func (repo *Experiments) getWithQuery(queryRow scannable) (*model.Experiment, error) {
	var exp model.Experiment

	err := queryRow.Scan(&exp.ID, &exp.Name, &exp.Description)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Error getting experiment from the DB: %v", err)
	}

	return &exp, nil
}

const selectExperimentsWhereIDSQL = `SELECT * FROM experiments WHERE id=$1`
const selectExperimentsSQL = `SELECT * FROM experiments`
const insertExperimentSQL = `INSERT INTO experiments (name, description) VALUES ($1, $2)`
const updateExperimentSQL = `UPDATE experiments SET name=$1, description=$2 WHERE id=$3`

// GetByID returns the Experiment with the given ID. If the Experiment does not
// exist, it returns nil, nil
func (repo *Experiments) GetByID(id int) (*model.Experiment, error) {
	return repo.getWithQuery(repo.db.QueryRow(selectExperimentsWhereIDSQL, id))
}

// GetAll returns all the Experiments
func (repo *Experiments) GetAll() ([]*model.Experiment, error) {
	rows, err := repo.db.Query(selectExperimentsSQL)
	if err != nil {
		return nil, fmt.Errorf("error getting experiments from the DB: %v", err)
	}
	defer rows.Close()

	results := make([]*model.Experiment, 0)

	for rows.Next() {
		exp, err := repo.getWithQuery(rows)
		if err != nil {
			return nil, fmt.Errorf("DB error: %v", err)
		}

		results = append(results, exp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("DB error: %v", err)
	}

	return results, nil
}

// Create experiment model in database. On success the assigned ID is set
func (repo *Experiments) Create(m *model.Experiment) error {
	r, err := repo.db.Exec(insertExperimentSQL, m.Name, m.Description)
	if err != nil {
		return err
	}
	newID, err := r.LastInsertId()
	if err != nil {
		return err
	}

	m.ID = int(newID)

	return nil
}

// Update experiment model in database
func (repo *Experiments) Update(m *model.Experiment) error {
	_, err := repo.db.Exec(updateExperimentSQL, m.Name, m.Description, m.ID)
	return err
}
