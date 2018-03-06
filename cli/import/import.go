/*
Tool to import pairs of files from an input sqlite database to another database.

Usage: import <path-to-sqlite.db> <destination-DSN>

Where DSN can be one of:
sqlite:///path/to/db.db
postgresql://[user[:password]@][netloc][:port][,...][/dbname]

The origin database is assumed to have the following table:
CREATE TABLE files (
	blob_id_a TEXT, repository_id_a TEXT, commit_hash_a TEXT, path_a TEXT, content_a TEXT, uast_a BLOB,
	blob_id_b TEXT, repository_id_b TEXT, commit_hash_b TEXT, path_b TEXT, content_b TEXT, uast_b BLOB,
	score DOUBLE PRECISION);
*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/src-d/code-annotation/server/dbutil"

	"github.com/jessevdk/go-flags"
)

// TODO: remove duplicated entries note if we decide to enforce the UNIQUE clause
const desc = `Imports pairs of files from the input database to the output database.
If the destination file does not exist, it will be created.

The Output argument must be one of:
sqlite:///path/to/db.db
postgresql://[user[:password]@][netloc][:port][,...][/dbname]

For a complete reference of the PostgreSQL connection string, see
https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING

The destination database does not need to be empty, new imported file pairs can
be added to previous imports.
Please note: if a file pair is identical to an existing one it will not be
detected. A new pair entry will be created with the same contents.`

var opts struct {
	Args struct {
		Input  string `description:"SQLite database filepath"`
		Output string `description:"SQLite or PostgreSQL Data Source Name"`
	} `positional-args:"yes" required:"yes"`
	ExperimentID int `long:"experiment-id" description:"Experiment ID to which files will be imported" required:"yes"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.LongDescription = desc

	if _, err := parser.Parse(); err != nil {
		if err, ok := err.(*flags.Error); ok {
			if err.Type == flags.ErrHelp {
				os.Exit(0)
			}

			fmt.Println()
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}

	originDB, err := dbutil.OpenSQLite(opts.Args.Input, true)
	if err != nil {
		log.Fatal(err)
	}
	defer originDB.Close()

	destDB, err := dbutil.Open(opts.Args.Output, false)
	if err != nil {
		log.Fatal(err)
	}
	defer destDB.Close()

	if err = dbutil.Bootstrap(destDB); err != nil {
		log.Fatal(err)
	}

	if err = dbutil.Initialize(destDB); err != nil {
		log.Fatal(err)
	}

	success, failures, err := dbutil.ImportFiles(originDB, destDB, dbutil.Options{}, opts.ExperimentID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Imported %v file pairs successfully\n", success)

	if failures > 0 {
		fmt.Printf("Failed to import %v file pairs\n", failures)
	}
}
