package main

import (
	"fmt"
	"log"

	"github.com/src-d/code-annotation/server/dbutil"
)

type vacuumCmd struct {
	commandDesc
	Args struct {
		InternalDBPath string `description:"filepath to the internal SQLite database"`
	} `positional-args:"yes" required:"yes"`
}

var vacuumOpts = vacuumCmd{
	commandDesc: commandDesc{
		name:      "vacuum",
		shortDesc: "Rebuilds the database to defragment it",
		longDesc:  "Rebuilds the sqlite://InternalDBPath database to eliminate free pages, compact table data...",
	},
}

// queries
const (
	vacuumQuery = "VACUUM"
)

func (c *vacuumCmd) Execute(args []string) error {
	internalDb, err := dbutil.Open(fmt.Sprintf(sqliteDSN, c.Args.InternalDBPath), false)
	if err != nil {
		log.Fatal(err)
	}

	defer internalDb.Close()

	log.Println("Running VACUUM process ...")

	_, err = internalDb.Exec(vacuumQuery)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("... VACUUM process finished")
	return nil
}
