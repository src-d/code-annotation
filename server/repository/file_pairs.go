package repository

import (
	"database/sql"
	"fmt"

	"github.com/src-d/code-annotation/server/model"
)

// FilePairs repository
type FilePairs struct {
	db *sql.DB
}

// NewFilePairs returns a new FilePairs repository
func NewFilePairs(db *sql.DB) *FilePairs {
	return &FilePairs{db: db}
}

// getWithQuery builds a FilePair from the given sql QueryRow. If the FilePair
// does not exist, it returns nil, nil
func (repo *FilePairs) getWithQuery(queryRow *sql.Row) (*model.FilePair, error) {
	var pair model.FilePair

	err := queryRow.Scan(&pair.ID,
		&pair.Left.BlobID, &pair.Left.RepositoryID, &pair.Left.CommitHash,
		&pair.Left.Path, &pair.Left.Content, &pair.Left.Hash,

		&pair.Right.BlobID, &pair.Right.RepositoryID, &pair.Right.CommitHash,
		&pair.Right.Path, &pair.Right.Content, &pair.Right.Hash,

		&pair.Score, &pair.Diff, &pair.ExperimentID)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("Error getting file pair from the DB: %v", err)
	default:
		return &pair, nil
	}
}

const selectFilePairsSQL = `SELECT * FROM file_pairs WHERE id=$1`

// GetByID returns the FilePair with the given ID. If the FilePair does not
// exist, it returns nil, nil
func (repo *FilePairs) GetByID(id int) (*model.FilePair, error) {
	return repo.getWithQuery(repo.db.QueryRow(selectFilePairsSQL, id))
}
