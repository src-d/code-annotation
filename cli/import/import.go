/*
Tool to import pairs of files from an input sqlite database to another database.

Usage: import <path-to-origin.db> <path-to-destination.db>

The origin database is assumed to have the following table:
CREATE TABLE files (name_a TEXT, name_b TEXT, content_a TEXT, content_b TEXT);
*/
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/src-d/code-annotation/server/dbutil"

	"github.com/jessevdk/go-flags"
)

// TODO: remove duplicated entries note if we decide to enforce the UNIQUE clause
const desc = `Imports pairs of files from the input database to the output database.
If the destination file does not exist, it will be created.

The destination database does not need to be empty, new imported file pairs can
be added to previous imports.
Please note: if a file pair is identical to an existing one it will not be
detected. A new pair entry will be created with the same contents.`

var opts struct {
	Args struct {
		Input  string `description:"SQLite database filepath"`
		Output string `description:"SQLite database filepath"`
	} `positional-args:"yes" required:"yes"`
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

	if _, err := os.Stat(opts.Args.Input); os.IsNotExist(err) {
		log.Fatalf("File %q does not exist", opts.Args.Input)
	}

	originDB, err := sql.Open("sqlite3", opts.Args.Input)
	if err != nil {
		log.Fatal(err)
	}
	defer originDB.Close()

	destDB, err := sql.Open("sqlite3", opts.Args.Output)
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

	success, failures, err := dbutil.ImportFiles(originDB, destDB, dbutil.Options{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Imported %v file pairs successfully\n", success)

	if failures > 0 {
		fmt.Printf("Failed to import %v file pairs\n", failures)
	}
}
