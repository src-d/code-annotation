package main

import (
	"fmt"
	"log"

	"github.com/src-d/code-annotation/server/dbutil"
)

type uastColsCmd struct {
	commandDesc
	Args struct {
		InternalDBPath string `description:"filepath to the internal SQLite database"`
	} `positional-args:"yes" required:"yes"`
}

var uastColsOpts = uastColsCmd{
	commandDesc: commandDesc{
		name:      "uast-add-cols",
		shortDesc: "Add uast BLOB columns",
		longDesc:  "Adds 'uast_a' and 'uast_b' BLOB columns to the sqlite://InternalDBPath database",
	},
}

// queries
const (
	leftFile  = "a"
	rightFile = "b"

	addUastQuery = "ALTER TABLE file_pairs ADD COLUMN uast_%s BLOB"
)

func (c *uastColsCmd) Execute(args []string) error {
	internalDb, err := dbutil.Open(fmt.Sprintf(sqliteDSN, c.Args.InternalDBPath), false)
	if err != nil {
		log.Fatal(err)
	}

	defer internalDb.Close()

	if err := addColumn(internalDb, leftFile); err != nil {
		log.Fatal(err)
	}

	if err := addColumn(internalDb, rightFile); err != nil {
		log.Fatal(err)
	}

	log.Println("New BLOB columns 'uast_a' and 'uast_b' were added")
	return nil
}

func addColumn(db dbutil.DB, side string) error {
	if _, err := db.Exec(fmt.Sprintf(addUastQuery, side)); err != nil {
		return err
	}

	return nil
}
