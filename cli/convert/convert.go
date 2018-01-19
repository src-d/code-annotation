// Tool to convert json files to an SQLite DB.

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

// The output DB will contain a single table with all the fields as columns,
// except for "bag":
// CREATE TABLE IF NOT EXISTS files (
//   blob_id_a TEXT, repository_id_a TEXT, commit_hash_a TEXT, path_a TEXT, content_a TEXT,
//   blob_id_b TEXT, repository_id_b TEXT, commit_hash_b TEXT, path_b TEXT, content_b TEXT,
//   score DOUBLE PRECISION)

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	_ "github.com/mattn/go-sqlite3"
)

const desc = `Converts a JSON file with pairs of files to the output database.
If the destination file does not exist, it will be created.`

var opts struct {
	Args struct {
		Input  string `description:"JSON file"`
		Output string `description:"SQLite database filepath"`
	} `positional-args:"yes" required:"yes"`
}

const (
	createTableSQL = `CREATE TABLE IF NOT EXISTS files (
		blob_id_a TEXT, repository_id_a TEXT, commit_hash_a TEXT, path_a TEXT, content_a TEXT,
		blob_id_b TEXT, repository_id_b TEXT, commit_hash_b TEXT, path_b TEXT, content_b TEXT,
		score DOUBLE PRECISION)`
	insertSQL = `INSERT INTO files VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
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

	source, err := ioutil.ReadFile(opts.Args.Input)
	if err != nil {
		log.Fatal(err)
	}

	destDB, err := sql.Open("sqlite3", opts.Args.Output)
	if err != nil {
		log.Fatal(err)
	}
	defer destDB.Close()

	_, err = destDB.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	var data [][3]interface{}
	if err := json.Unmarshal(source, &data); err != nil {
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

	for _, pair := range data {
		fileA := pair[0].(map[string]interface{})
		fileB := pair[1].(map[string]interface{})

		res, err := insert.Exec(
			fileA["blob_id"], fileA["repository_id"], fileA["commit_hash"], fileA["path"], fileA["content"],
			fileB["blob_id"], fileB["repository_id"], fileB["commit_hash"], fileB["path"], fileB["content"],
			pair[2].(float64))

		if err != nil {
			failures++
			log.Println(err)
		}

		rowsAffected, _ := res.RowsAffected()
		success += rowsAffected
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit the transaction:\n%v", err)
	}

	fmt.Printf("Converted %v file pairs successfully\n", success)

	if failures > 0 {
		fmt.Printf("Failed to convert %v file pairs\n", failures)
	}
}
