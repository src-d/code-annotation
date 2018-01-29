package dbutil

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pmezard/go-difflib/difflib"
)

type driver int

const (
	none driver = iota
	sqlite
	postgres
)

// DB groups a sql.DB and the driver used to initialize it
type DB struct {
	sqlDB  *sql.DB
	driver driver
}

// Close closes the database, releasing any open resources.
func (db *DB) Close() error {
	return db.sqlDB.Close()
}

const (
	incrementTypePlaceholder = "<INCREMENT_TYPE>"
	sqliteIncrementType      = "INTEGER"
	posgresIncrementType     = "SERIAL"

	createUsers = `CREATE TABLE IF NOT EXISTS users (
			id <INCREMENT_TYPE>, github_username TEXT, auth TEXT, role INTEGER,
			PRIMARY KEY (id))`
	createExperiments = `CREATE TABLE IF NOT EXISTS experiments (
			id <INCREMENT_TYPE>, name TEXT UNIQUE, description TEXT,
			PRIMARY KEY (id))`
	// TODO: consider a unique constrain to avoid importing identical pairs
	createFilePairs = `CREATE TABLE IF NOT EXISTS file_pairs (
		id <INCREMENT_TYPE>,
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

const (
	defaultExperimentID = 1

	insertExperiments = `INSERT INTO experiments
		(id, name, description)
		VALUES ($1, 'default', 'Default experiment')`

	alterExperimentsSequence = `ALTER SEQUENCE experiments_id_seq RESTART WITH 2`
)

const selectFiles = `SELECT * FROM files`

const insertFilePairs = `INSERT INTO file_pairs (
		blob_id_a, repository_id_a, commit_hash_a, path_a, content_a, hash_a,
		blob_id_b, repository_id_b, commit_hash_b, path_b, content_b, hash_b,
		score, diff, experiment_id ) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

var (
	sqliteReg = regexp.MustCompile(`^sqlite://(.+)$`)
	psReg     = regexp.MustCompile(`^postgres(ql)?:.+$`)
)

// OpenSQLite calls Open, but using a path to an SQLite file. For sintactic
// sugar a sqlite://path string is also accepted
func OpenSQLite(filepath string, checkExisting bool) (DB, error) {
	if psReg.MatchString(filepath) {
		return DB{nil, none}, fmt.Errorf(
			"Invalid PostgreSQL connection string %q, a path to an SQLite file was expected",
			filepath)
	}

	if !sqliteReg.MatchString(filepath) {
		filepath = `sqlite://` + filepath
	}

	return Open(filepath, checkExisting)
}

// Open returns a DB from the connection string.
// With checkExisting it will fail if the DB does not exist
func Open(connection string, checkExisting bool) (DB, error) {
	if conn := sqliteReg.FindStringSubmatch(connection); conn != nil {
		if checkExisting {
			if _, err := os.Stat(conn[1]); os.IsNotExist(err) {
				return DB{nil, none}, fmt.Errorf("File %q does not exist", conn[1])
			}
		}

		db, err := sql.Open("sqlite3", conn[1])
		return DB{db, sqlite}, err
	}

	if psReg.MatchString(connection) {
		db, err := sql.Open("postgres", connection)
		return DB{db, postgres}, err
	}

	return DB{nil, none}, fmt.Errorf(`Connection string %q is not valid. It must be on of
sqlite:///path/to/db.db
postgresql://[user[:password]@][netloc][:port][,...][/dbname]`, connection)
}

// Bootstrap creates the necessary tables for the output DB. It is safe to call on a
// DB that is already bootstrapped.
func Bootstrap(db DB) error {
	tables := []string{createUsers, createExperiments,
		createFilePairs, createAssignments}

	var colType string

	switch db.driver {
	case sqlite:
		colType = sqliteIncrementType
	case postgres:
		colType = posgresIncrementType
	default:
		return fmt.Errorf("Unknown driver type")
	}

	for _, table := range tables {
		cmd := strings.Replace(table, incrementTypePlaceholder, colType, 1)

		if _, err := db.sqlDB.Exec(cmd); err != nil {
			return err
		}
	}

	return nil
}

// Initialize populates the DB with default values. It is safe to call on a
// DB that is already initialized
func Initialize(db DB) error {
	_, err := db.sqlDB.Exec(insertExperiments, defaultExperimentID)
	if db.driver == postgres && err == nil {
		db.sqlDB.Exec(alterExperimentsSequence)
	}

	// Errors are ignored to allow initialization over an existing DB
	return nil
}

// Options for the ImportFiles and Copy methods.
// Logger is optional, if it is not provided the default stderr will be used.
type Options struct {
	Logger *log.Logger
}

func (opts *Options) getLogger() *log.Logger {
	if opts.Logger != nil {
		return opts.Logger
	}

	return log.New(os.Stderr, "", log.LstdFlags) // Default log to stderr

}

// ImportFiles imports pairs of files from the origin to the destination DB.
// It copies the contents and processes the needed data (md5 hash, diff)
func ImportFiles(originDB DB, destDB DB, opts Options) (success, failures int64, e error) {

	logger := opts.getLogger()

	rows, err := originDB.sqlDB.Query(selectFiles)
	if err != nil {
		return 0, 0, err
	}

	tx, err := destDB.sqlDB.Begin()
	if err != nil {
		return 0, 0, err
	}

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
			logger.Printf("Failed to insert row\nerror: %v\n", err)
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
