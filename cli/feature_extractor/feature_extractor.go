// Tool to extract features from json files to an SQLite DB.

// Beware, this is an internal tool. Do not expect production-ready
// documentation, maintenance, or tests.

// This is how the json file looks like:
// [
//   [
//     {
//       "blob_id": "3a6e6[..]196",
//       "repository_id": "github.com/[..]git",
//       "commit_hash": "92[..]9d5",
//       "bag": {
//         "r.[..]": 0.45368270469072236,
//         "r.[..]": 1.8373093877436146,
// 			...
//       }
//       "path": "/src/[..].java",
//       "content": "file contents"
//     },
//     {
//       "blob_id": "3a6e6[..]196",
//       "repository_id": "github.com/[..]git",
//       "commit_hash": "92[..]9d5",
//       "bag": {
//         "r.[..]": 0.45368270469072236,
//         "r.[..]": 1.8373093877436146,
// 			...
//       }
//       "path": "/src/[..].java",
//       "content": "file contents"
//     },
//     0.9512810301340767
//   ],
//   ...
// ]

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	_ "github.com/mattn/go-sqlite3"
	"github.com/src-d/code-annotation/server/dbutil"
)

const desc = `Extracts features from a JSON file with pairs of files to the output database.
The destination db must exist.

The Output argument must be one of:
sqlite:///path/to/db.db
postgresql://[user[:password]@][netloc][:port][,...][/dbname]

For a complete reference of the PostgreSQL connection string, see
https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING`

var opts struct {
	Args struct {
		Input  string `description:"JSON file"`
		Output string `description:"SQLite database filepath or PostgreSQL Data Source Name"`
	} `positional-args:"yes" required:"yes"`
}

const (
	createFeaturesSQL = `CREATE TABLE IF NOT EXISTS features (
		blob_id TEXT,
		name TEXT, weight REAL,
		PRIMARY KEY (blob_id, name))`
	insertSQL = `INSERT INTO features VALUES ($1, $2, $3)`
)

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

	destDB, err := dbutil.Open(opts.Args.Output, true)
	if err != nil {
		log.Fatal(err)
	}
	defer destDB.Close()

	source, err := os.OpenFile(opts.Args.Input, os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	_, err = destDB.Exec(createFeaturesSQL)
	if err != nil {
		log.Fatal(err)
	}

	var data [][3]interface{}
	if err := json.NewDecoder(source).Decode(&data); err != nil {
		log.Fatal(err)
	}

	tx, err := destDB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	insert, err := tx.Prepare(insertSQL)
	if err != nil {
		log.Fatal(err)
	}

	var success, failures int64

	// there are full duplicates in json (checked with map & reflect.DeepEqual)
	// instead of printing many errors like UNIQUE constraint failed: features.blob_id, features.name
	// I handle them here
	blobIDs := make(map[string]bool)

	processFile := func(f map[string]interface{}) {
		blobID := f["blob_id"].(string)
		if _, ok := blobIDs[blobID]; ok {
			return
		}
		blobIDs[blobID] = true

		features := f["bag"].(map[string]interface{})
		for name, weight := range features {
			res, err := insert.Exec(blobID, name, weight.(float64))

			if err != nil {
				failures++
				log.Println(err)
				continue
			}

			rowsAffected, _ := res.RowsAffected()
			success += rowsAffected
		}
	}

	for _, pair := range data {
		processFile(pair[0].(map[string]interface{}))
		processFile(pair[1].(map[string]interface{}))
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit the transaction:\n%v", err)
	}

	fmt.Printf("Extracted and saved %v features successfully\n", success)

	if failures > 0 {
		fmt.Printf("Failed to extract and save %v features\n", failures)
	}
}
