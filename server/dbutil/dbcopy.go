package dbutil

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	tablePlaceholder = `<TABLE>`

	dumpAllSQL       = `SELECT * FROM <TABLE>`
	bulkInsertSQL    = `INSERT INTO <TABLE> VALUES `
	maxIDSQL         = `SELECT MAX(id) FROM <TABLE>`
	alterSequenceSQL = `ALTER SEQUENCE <TABLE>_id_seq RESTART WITH $1`
)

var tables = []string{"users", "experiments", "file_pairs", "assignments"}

// Copy dumps the contents of the origin DB into the destination DB. The
// destination DB should be bootstrapped, but empty
func Copy(originDB DB, destDB DB, opts Options) error {

	logger := opts.getLogger()

	tx, err := destDB.Begin()
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			logger.Println("The transaction will be rolled back")
			tx.Rollback()
		}
	}()

	for _, table := range tables {
		// SELECT * FROM <TABLE>
		selectCmd := strings.Replace(dumpAllSQL, tablePlaceholder, table, 1)

		rows, err := originDB.Query(selectCmd)
		if err != nil {
			return err
		}
		defer rows.Close()

		columnNames, _ := rows.Columns()
		nColumns := len(columnNames)

		// INSERT INTO <TABLE> VALUES ($1,$2...)
		insert, err := tx.Prepare(insertCmd(table, nColumns))
		if err != nil {
			return err
		}

		// Generic arguments to read and write the values
		columnValsPtr := genericVals(nColumns)

		var success int64

		for rows.Next() {
			if err := rows.Scan(columnValsPtr...); err != nil {
				return err
			}

			res, err := insert.Exec(columnValsPtr...)
			if err != nil {
				return err
			}

			rowsAffected, _ := res.RowsAffected()
			success += rowsAffected
		}

		logger.Printf("Inserted %v rows into table %v\n", success, table)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	committed = true

	fixSequences(destDB, logger)

	return nil
}

// nArgs returns a string containing INSERT INTO <TABLE> VALUES ($1,$2...)
func insertCmd(table string, n int) string {
	cmd := strings.Replace(bulkInsertSQL, tablePlaceholder, table, 1)

	var nArgs []string

	for i := 1; i <= n; i++ {
		nArgs = append(nArgs, "$"+strconv.Itoa(i))
	}

	return cmd + "(" + strings.Join(nArgs, ",") + ")"
}

// genericVals returns a slice of interface{}, each one a pointer to a NullString
func genericVals(nColumns int) []interface{} {
	columnVals := make([]sql.NullString, nColumns)
	columnValsPtr := make([]interface{}, nColumns)

	for i := range columnVals {
		columnValsPtr[i] = &columnVals[i]
	}

	return columnValsPtr
}

// fixSequences fixes the PostgreSQL sequences setting them to the max id found
// This should not be needed for readonly DBs, but his way the DB is ready in
// case new rows are added
func fixSequences(db DB, logger *log.Logger) {
	if db.driver != postgres {
		return
	}

	for _, table := range tables {
		selectCmd := strings.Replace(maxIDSQL, tablePlaceholder, table, 1)

		var maxID sql.NullInt64
		err := db.QueryRow(selectCmd).Scan(&maxID)

		if err != nil {
			// With 0 rows there is no MAX(id)
			continue
		}

		if maxID.Valid {
			alterCmd := strings.Replace(alterSequenceSQL, tablePlaceholder, table, 1)

			newMax := fmt.Sprintf("%v", maxID.Int64+1)

			// for some reason Exec() fails to substitute $1 with the argument
			//_, err := db.Exec(alterCmd, newMax)

			alterCmd = strings.Replace(alterCmd, "$1", newMax, 1)
			_, err := db.Exec(alterCmd)

			if err != nil {
				logger.Printf("Error while executing %q:\n%v\n", alterCmd, err)
			}
		}
	}
}
