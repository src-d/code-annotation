package dbutil

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pmezard/go-difflib/difflib"
)

const defaultExperimentID = 1

const (
	createUsers = `CREATE TABLE IF NOT EXISTS users (
			id INTEGER, github_username TEXT, auth TEXT, role INTEGER,
			PRIMARY KEY (id))`
	createExperiments = `CREATE TABLE IF NOT EXISTS experiments (
			id INTEGER, name TEXT UNIQUE, description TEXT,
			PRIMARY KEY (id))`
	// TODO: consider a unique constrain to avoid importing identical pairs
	createFilePairs = `CREATE TABLE IF NOT EXISTS file_pairs (
		id INTEGER,
		blob_id_a TEXT, repository_id_a TEXT, commit_hash_a TEXT, path_a TEXT, content_a TEXT, hash_a TEXT,
		blob_id_b TEXT, repository_id_b TEXT, commit_hash_b TEXT, path_b TEXT, content_b TEXT, hash_b TEXT,
		score DOUBLE PRECISION, diff TEXT, experiment_id INTEGER,
		PRIMARY KEY (id),
		FOREIGN KEY(experiment_id) REFERENCES experiments(id))`
	createAssignments = `CREATE TABLE IF NOT EXISTS assignments (
			user_id INTEGER, pair_id INTEGER, experiment_id INTEGER,
			answer INTEGER, duration INTEGER,
			PRIMARY KEY (user_id, pair_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (pair_id) REFERENCES file_pairs(id),
			FOREIGN KEY (experiment_id) REFERENCES experiments(id))`
)

const insertExperiments = `INSERT OR IGNORE INTO experiments
	(id, name, description)
	VALUES ($1, 'default', 'Default experiment')`

const selectFiles = `SELECT * FROM files`

const insertFilePairs = `INSERT INTO file_pairs (
		blob_id_a, repository_id_a, commit_hash_a, path_a, content_a, hash_a,
		blob_id_b, repository_id_b, commit_hash_b, path_b, content_b, hash_b,
		score, diff, experiment_id ) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

// Bootstrap creates the necessary tables for the output DB. It is safe to call on a
// DB that is already bootstrapped.
func Bootstrap(db *sql.DB) error {
	tables := []string{createUsers, createExperiments,
		createFilePairs, createAssignments}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}

	return nil
}

// Initialize populates the DB with default values. It is safe to call on a
// DB that is already initialized
func Initialize(db *sql.DB) error {
	_, err := db.Exec(insertExperiments, defaultExperimentID)
	return err
}

// Options for the ImportFiles method. Logger is optional, if it is not provided
// the default stderr will be used.
type Options struct {
	Logger *log.Logger
}

// ImportFiles imports pairs of files from the origin to the destination DB.
// It copies the contents and processes the needed data (md5 hash, diff)
func ImportFiles(originDB, destDB *sql.DB, opts Options) (success, failures int64, e error) {
	var logger *log.Logger
	if opts.Logger != nil {
		logger = opts.Logger
	} else {
		logger = log.New(os.Stderr, "", log.LstdFlags) // Default log to stderr
	}

	rows, err := originDB.Query(selectFiles)
	if err != nil {
		return 0, 0, err
	}

	tx, err := destDB.Begin()

	insert, err := tx.Prepare(insertFilePairs)

	if err != nil {
		return 0, 0, err
	}

	for rows.Next() {
		var blobIDA, repositoryIDA, commitHashA, pathA, contentA,
			blobIDB, repositoryIDB, commitHashB, pathB, contentB string
		var score float64

		err := rows.Scan(
			&blobIDA, &repositoryIDA, &commitHashA, &pathA, &contentA,
			&blobIDB, &repositoryIDB, &commitHashB, &pathB, &contentB,
			&score)

		if err != nil {
			logger.Printf("Failed to read row from origin DB\nerror: %v\n", err)
			failures++
			continue
		}

		diffText, err := diff(pathA, pathB, contentA, contentB)
		if err != nil {
			logger.Printf(
				"Failed to create diff for files:\n - %q\n - %q\nerror: %v\n",
				pathA, pathB, err)
			failures++
			continue
		}

		res, err := insert.Exec(
			blobIDA, repositoryIDA, commitHashA, pathA, contentA, md5hash(contentA),
			blobIDB, repositoryIDB, commitHashB, pathB, contentB, md5hash(contentB),
			score,
			diffText,
			defaultExperimentID)

		if err != nil {
			logger.Println(err)
			failures++
			continue
		}

		rowsAffected, _ := res.RowsAffected()
		success += rowsAffected
	}

	if err := tx.Commit(); err != nil {
		return 0, success + failures, err
	}

	return success, failures, rows.Err()
}

func md5hash(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}

func diff(nameA, nameB, contentA, contentB string) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(contentA),
		B:        difflib.SplitLines(contentB),
		FromFile: nameA,
		ToFile:   nameB,
		Context:  3,
	}

	return difflib.GetUnifiedDiffString(diff)
}
