package dbutil

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	// loads the driver
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Driver used
type Driver int

// Available drivers
const (
	None Driver = iota
	Sqlite
	Postgres
)

// DB groups a sql.DB and the driver used to initialize it
type DB struct {
	*sql.DB
	Driver Driver
}

// SQLDB returns the *sql.DB
func (db *DB) SQLDB() *sql.DB {
	return db.DB
}

const (
	incrementTypePlaceholder = "<INCREMENT_TYPE>"
	sqliteIncrementType      = "INTEGER"
	posgresIncrementType     = "SERIAL"

	blobTypePlaceholder = "<BLOB_TYPE>"
	sqliteBlobType      = "BLOB"
	posgresBlobType     = "BYTEA"

	createUsers = `CREATE TABLE IF NOT EXISTS users (
			id <INCREMENT_TYPE>, login TEXT UNIQUE, username TEXT, avatar_url TEXT, role TEXT,
			PRIMARY KEY (id))`
	createExperiments = `CREATE TABLE IF NOT EXISTS experiments (
			id <INCREMENT_TYPE>, name TEXT UNIQUE, description TEXT,
			PRIMARY KEY (id))`
	// TODO: consider a unique constrain to avoid importing identical pairs
	createFilePairs = `CREATE TABLE IF NOT EXISTS file_pairs (
		id <INCREMENT_TYPE>,
		blob_id_a TEXT, repository_id_a TEXT, commit_hash_a TEXT, path_a TEXT, content_a TEXT, hash_a TEXT,
		blob_id_b TEXT, repository_id_b TEXT, commit_hash_b TEXT, path_b TEXT, content_b TEXT, hash_b TEXT,
		score DOUBLE PRECISION, experiment_id INTEGER,
		uast_a <BLOB_TYPE>, uast_b <BLOB_TYPE>,
		PRIMARY KEY (id),
		FOREIGN KEY(experiment_id) REFERENCES experiments(id))`
	createAssignments = `CREATE TABLE IF NOT EXISTS assignments (
			id <INCREMENT_TYPE>,
			user_id INTEGER, pair_id INTEGER, experiment_id INTEGER,
			answer TEXT, duration INTEGER,
			PRIMARY KEY (id),
			UNIQUE (user_id, pair_id, experiment_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (pair_id) REFERENCES file_pairs(id),
			FOREIGN KEY (experiment_id) REFERENCES experiments(id))`
	createFeatures = `CREATE TABLE IF NOT EXISTS features (
		blob_id TEXT,
		name TEXT, weight REAL,
		PRIMARY KEY (blob_id, name))`
)

const (
	defaultExperimentID = 1

	insertExperiments = `INSERT INTO experiments
		(id, name, description)
		VALUES ($1, 'default', 'Default experiment')`

	alterExperimentsSequence = `ALTER SEQUENCE experiments_id_seq RESTART WITH 2`
)

const selectExperiment = `SELECT * FROM experiments WHERE id = $1`

const selectFiles = `SELECT
	blob_id_a, repository_id_a, commit_hash_a, path_a, content_a, uast_a,
	blob_id_b, repository_id_b, commit_hash_b, path_b, content_b, uast_b,
	score FROM files`

const insertFilePairs = `INSERT INTO file_pairs (
		blob_id_a, repository_id_a, commit_hash_a, path_a, content_a, hash_a, uast_a,
		blob_id_b, repository_id_b, commit_hash_b, path_b, content_b, hash_b, uast_b,
		score, experiment_id ) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

var (
	sqliteReg = regexp.MustCompile(`^sqlite://(.+)$`)
	psReg     = regexp.MustCompile(`^postgres(ql)?:.+$`)
)

// OpenSQLite calls Open, but using a path to an SQLite file. For sintactic
// sugar a sqlite://path string is also accepted
func OpenSQLite(filepath string, checkExisting bool) (DB, error) {
	if psReg.MatchString(filepath) {
		return DB{nil, None}, fmt.Errorf(
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
				return DB{nil, None}, fmt.Errorf("File %q does not exist", conn[1])
			}
		}

		db, err := sql.Open("sqlite3", conn[1])
		return DB{db, Sqlite}, err
	}

	if psReg.MatchString(connection) {
		db, err := sql.Open("postgres", connection)
		return DB{db, Postgres}, err
	}

	return DB{nil, None}, fmt.Errorf(`Connection string %q is not valid. It must be on of
sqlite:///path/to/db.db
postgresql://[user[:password]@][netloc][:port][,...][/dbname]`, connection)
}

// Bootstrap creates the necessary tables for the output DB. It is safe to call on a
// DB that is already bootstrapped.
func Bootstrap(db DB) error {
	tables := []string{createUsers, createExperiments,
		createFilePairs, createAssignments, createFeatures}

	var colType string
	var blobType string

	switch db.Driver {
	case Sqlite:
		colType = sqliteIncrementType
		blobType = sqliteBlobType
	case Postgres:
		colType = posgresIncrementType
		blobType = posgresBlobType
	default:
		return fmt.Errorf("Unknown driver type")
	}

	for _, table := range tables {
		cmd := strings.Replace(table, incrementTypePlaceholder, colType, 1)
		cmd = strings.Replace(cmd, blobTypePlaceholder, blobType, -1)

		if _, err := db.Exec(cmd); err != nil {
			return err
		}
	}

	return nil
}

// Initialize populates the DB with default values. It is safe to call on a
// DB that is already initialized
func Initialize(db DB) error {
	_, err := db.Exec(insertExperiments, defaultExperimentID)
	if db.Driver == Postgres && err == nil {
		db.Exec(alterExperimentsSequence)
	}

	// Errors are ignored to allow initialization over an existing DB
	return nil
}

// Options for the ImportFiles and Copy methods.
// Logger is optional, if it is not provided the default stderr will be used.
type Options struct {
	Logger logrus.FieldLogger
}

func (opts *Options) getLogger() logrus.FieldLogger {
	if opts.Logger != nil {
		return opts.Logger
	}

	return logrus.StandardLogger()
}

// ImportFiles imports pairs of files from the origin to the destination DB.
// It copies the contents and processes the needed data (md5 hash)
func ImportFiles(originDB DB, destDB DB, opts Options, experimentID int) (success, failures int64, e error) {

	logger := opts.getLogger()

	rows, err := destDB.Query(selectExperiment, experimentID)
	if err != nil {
		return 0, 0, err
	}

	if !rows.Next() {
		rows.Close()
		return 0, 0, fmt.Errorf("Experiment with id %d doesn't exist", experimentID)
	}
	rows.Close()

	rows, err = originDB.Query(selectFiles)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	tx, err := destDB.Begin()
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
		var uastA, uastB []byte
		var score float64

		err := rows.Scan(
			&blobIDA, &repositoryIDA, &commitHashA, &pathA, &contentA, &uastA,
			&blobIDB, &repositoryIDB, &commitHashB, &pathB, &contentB, &uastB,
			&score)

		if err != nil {
			logger.Printf("Failed to read row from origin DB\nerror: %v\n", err)
			failures++
			continue
		}

		res, err := insert.Exec(
			blobIDA, repositoryIDA, commitHashA, pathA, contentA, md5hash(contentA), uastA,
			blobIDB, repositoryIDB, commitHashB, pathB, contentB, md5hash(contentB), uastB,
			score,
			experimentID)

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
