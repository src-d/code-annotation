package repository

import (
	"database/sql"
	"fmt"

	"github.com/src-d/code-annotation/server/model"
)

// Features repository
type Features struct {
	db *sql.DB
}

// NewFeatures returns a new Features repository
func NewFeatures(db *sql.DB) *Features {
	return &Features{db: db}
}

const selectFeaturesSQL = `SELECT name, weight FROM features WHERE blob_id=$1`

// GetAll returns a list of all features for blobID
func (repo *Features) GetAll(blobID string) ([]*model.Feature, error) {
	rows, err := repo.db.Query(selectFeaturesSQL, blobID)
	if err != nil {
		return nil, fmt.Errorf("Error getting features from the DB: %v", err)
	}
	defer rows.Close()

	results := make([]*model.Feature, 0)

	for rows.Next() {
		var f model.Feature
		rows.Scan(&f.Name, &f.Weight)

		results = append(results, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("DB error: %v", err)
	}

	return results, nil
}
