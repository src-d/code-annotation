package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/src-d/code-annotation/server/dbutil"
)

type diffRmColCmd struct {
	commandDesc
	Args struct {
		InternalDBPath string `description:"filepath to the internal SQLite database"`
	} `positional-args:"yes" required:"yes"`
}

var diffRmColOpts = diffRmColCmd{
	commandDesc: commandDesc{
		name:      "diff-rm-col",
		shortDesc: "Remove diff column",
		longDesc:  "Remove diff column from sqlite://InternalDBPath database",
	},
}

// queries
const (
	fileFairsTable = "file_pairs"
	tmpTable       = "file_pairs_tmp"

	cols = `id,
		blob_id_a, repository_id_a, commit_hash_a, path_a, content_a, hash_a, uast_a,
		blob_id_b, repository_id_b, commit_hash_b, path_b, content_b, hash_b, uast_b,
		score, experiment_id`

	createTmpTableQuery = `CREATE TABLE IF NOT EXISTS ` + tmpTable + ` (
		id INTEGER,
		blob_id_a TEXT, repository_id_a TEXT, commit_hash_a TEXT, path_a TEXT, content_a TEXT, hash_a TEXT,
		blob_id_b TEXT, repository_id_b TEXT, commit_hash_b TEXT, path_b TEXT, content_b TEXT, hash_b TEXT,
		score DOUBLE PRECISION, experiment_id INTEGER,
		uast_a BLOB, uast_b BLOB,
		PRIMARY KEY (id),
		FOREIGN KEY(experiment_id) REFERENCES experiments(id))`

	fillTmpTableQuery   = "INSERT INTO " + tmpTable + "(" + cols + ") SELECT " + cols + " FROM " + fileFairsTable
	disableIndexQuery   = "PRAGMA foreign_keys=OFF"
	dropOldTableQuery   = "DROP TABLE " + fileFairsTable
	renameTmpTableQuery = "ALTER TABLE " + tmpTable + " RENAME TO " + fileFairsTable
	checkIndexQuery     = "PRAGMA foreign_key_check"
	enableIndexQuery    = "PRAGMA foreign_keys=ON"
)

func (c *diffRmColCmd) Execute(args []string) error {
	internalDb, err := dbutil.Open(fmt.Sprintf(sqliteDSN, c.Args.InternalDBPath), false)
	if err != nil {
		log.Fatal(err)
	}

	defer internalDb.Close()

	queries := []string{
		createTmpTableQuery,
		fillTmpTableQuery,
		disableIndexQuery,
		dropOldTableQuery,
		renameTmpTableQuery,
		enableIndexQuery,
	}

	if err := execQueries(internalDb, queries); err != nil {
		log.Fatal(err)
	}

	log.Println("Deleted 'diff' column")
	return nil
}

func execQueries(db dbutil.DB, queries []string) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("Error running migration; %s \nThe rollback failed; %s", err, rollbackErr)
			}
		}
	}()

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			return err
		}
	}

	if err := ensureForeignKeys(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func ensureForeignKeys(tx queryer) error {
	rows, err := tx.Query(checkIndexQuery)
	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("Foreign key constraints were violated")
	}

	return nil
}

type queryer interface {
	Query(string, ...interface{}) (*sql.Rows, error)
}
