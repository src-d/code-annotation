package repository

import (
	"database/sql"
	"fmt"

	"github.com/src-d/code-annotation/server/model"
)

// Assignments repository
type Assignments struct {
	db *sql.DB
}

// NewAssignments returns a new Assignments repository
func NewAssignments(db *sql.DB) *Assignments {
	return &Assignments{db: db}
}

const (
	insertAssignmentsSQL             = `INSERT INTO assignments (user_id, pair_id, experiment_id, answer, duration) VALUES ($1, $2, $3, $4, $5)`
	selectIDFilePairsSQL             = `SELECT id FROM file_pairs WHERE experiment_id=$1 AND id NOT IN (SELECT pair_id FROM assignments WHERE experiment_id=$1 AND user_id=$2)`
	selectAssignmentsSQL             = `SELECT * FROM assignments WHERE user_id=$1 AND experiment_id=$2`
	selectAssignmentsWhereExpPairSQL = `SELECT * FROM assignments WHERE experiment_id=$1 AND pair_id=$2`
	updateAssignmentsSQL             = `UPDATE assignments SET answer='%v', duration=%v WHERE id=%v`
	countPendingIDsSQL               = `SELECT count(id) FROM file_pairs WHERE experiment_id=$1 AND id NOT IN (SELECT pair_id FROM assignments WHERE experiment_id=$1 AND user_id=$2)`
)

// IsInitialized returns true if the assignments are initialized for the given
// user and experiment IDs. If it's false, Initialize should be called
func (repo *Assignments) IsInitialized(userID, experimentID int) (bool, error) {
	row := repo.db.QueryRow(countPendingIDsSQL, experimentID, userID)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	return count == 0, nil
}

// Initialize builds the assignments for the given user and experiment IDs
func (repo *Assignments) Initialize(userID int, experimentID int) (int, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return 0, err
	}

	insert, err := tx.Prepare(insertAssignmentsSQL)
	if err != nil {
		return 0, fmt.Errorf("DB error: %v", err)
	}

	rows, err := repo.db.Query(selectIDFilePairsSQL, experimentID, userID)
	if err != nil {
		return 0, fmt.Errorf("Error getting file_pairs from the DB: %v", err)
	}
	defer rows.Close()

	duration := 0
	created := 0
	for rows.Next() {
		var pairID int
		rows.Scan(&pairID)
		_, err := insert.Exec(userID, pairID, experimentID, nil, duration)
		if err != nil {
			return 0, fmt.Errorf("DB error: %v", err)
		}

		created++
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("DB error: %v", err)
	}

	return created, nil
}

// getWithQuery builds a Assignment from the given sql QueryRow. If the
// Assignment does not exist, it returns nil, nil
func (repo *Assignments) getWithQuery(queryRow scannable) (*model.Assignment, error) {
	var as model.Assignment

	err := queryRow.Scan(&as.ID, &as.UserID, &as.PairID, &as.ExperimentID,
		&as.Answer, &as.Duration)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("Error getting assignment from the DB: %v", err)
	default:
		return &as, nil
	}
}

// GetByID returns the Assignment with the given ID. If the Assignment does not
// exist, it returns nil, nil
func (repo *Assignments) GetByID(id int) (*model.Assignment, error) {
	return repo.getWithQuery(
		repo.db.QueryRow("SELECT * FROM assignments WHERE id=$1", id))
}

func (repo *Assignments) getAssignmentsWithQuery(query string, args ...interface{}) ([]*model.Assignment, error) {
	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting assignments from the DB: %v", err)
	}
	defer rows.Close()

	results := make([]*model.Assignment, 0)

	for rows.Next() {
		as, err := repo.getWithQuery(rows)
		if err != nil {
			return nil, fmt.Errorf("DB error: %v", err)
		}

		results = append(results, as)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("DB error: %v", err)
	}

	return results, nil
}

// GetAll returns all the Assignments for the given user and experiment IDs
func (repo *Assignments) GetAll(userID, experimentID int) ([]*model.Assignment, error) {
	return repo.getAssignmentsWithQuery(selectAssignmentsSQL, userID, experimentID)
}

// GetByExperimentPair returns all the Assignments for the given experiment and pair IDs
func (repo *Assignments) GetByExperimentPair(experimentID, filePairID int) ([]*model.Assignment, error) {
	return repo.getAssignmentsWithQuery(
		selectAssignmentsWhereExpPairSQL, experimentID, filePairID)
}

// Update updates the Assignment identified by the given user and pair IDs,
// with the given answer and duration
func (repo *Assignments) Update(assignmentID int, answer string, duration int) error {
	if _, ok := model.Answers[answer]; !ok {
		return fmt.Errorf("Wrong answer provided: '%s'", answer)
	}

	cmd := fmt.Sprintf(updateAssignmentsSQL, answer, duration, assignmentID)

	_, err := repo.db.Exec(cmd)

	return err
}
