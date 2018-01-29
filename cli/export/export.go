/*
Tool to export annotation results from the internal DB to an output DB.
It does a simple copy&paste of the internal DB. The annotation results are
stored in the assignments table.

Usage: import <origin-DSN> <path-to-sqlite.db>

Where DSN can be one of:
sqlite:///path/to/db.db
postgresql://[user[:password]@][netloc][:port][,...][/dbname]
*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/src-d/code-annotation/server/dbutil"

	"github.com/jessevdk/go-flags"
)

const desc = `Exports annotation results from the internal input database to a new output
database. The destination database must be empty.

The Input argument must be one of:
sqlite:///path/to/db.db
postgresql://[user[:password]@][netloc][:port][,...][/dbname]

For a complete reference of the PostgreSQL connection string, see
https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING`

var opts struct {
	Args struct {
		Input  string `description:"SQLite or PostgreSQL Data Source Name"`
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

	originDB, err := dbutil.Open(opts.Args.Input, true)
	if err != nil {
		log.Fatal(err)
	}
	defer originDB.Close()

	destDB, err := dbutil.OpenSQLite(opts.Args.Output, false)
	if err != nil {
		log.Fatal(err)
	}
	defer destDB.Close()

	if err = dbutil.Bootstrap(destDB); err != nil {
		log.Fatal(err)
	}

	if err := dbutil.Copy(originDB, destDB, dbutil.Options{}); err != nil {
		log.Fatal(err)
	}
}
