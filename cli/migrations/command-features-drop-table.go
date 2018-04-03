package main

import (
	"fmt"
	"log"

	"github.com/src-d/code-annotation/server/dbutil"
)

type featureDropCmd struct {
	commandDesc
	Args struct {
		InternalDBPath string `description:"filepath to the internal SQLite database"`
	} `positional-args:"yes" required:"yes"`
}

var featureDropOpts = featureDropCmd{
	commandDesc: commandDesc{
		name:      "features-drop-table",
		shortDesc: "Drop Features table",
		longDesc:  "Removes the Features table from the sqlite://InternalDBPath database",
	},
}

// queries
const (
	dropFeaturesTableQuery = "DROP TABLE features"
)

func (c *featureDropCmd) Execute(args []string) error {
	internalDb, err := dbutil.Open(fmt.Sprintf(sqliteDSN, c.Args.InternalDBPath), false)
	if err != nil {
		log.Fatal(err)
	}

	defer internalDb.Close()

	if err := dropFeaturesTable(internalDb); err != nil {
		log.Fatal(err)
	}

	log.Println("Features table was deleted")
	return nil
}

func dropFeaturesTable(db dbutil.DB) error {
	if _, err := db.Exec(dropFeaturesTableQuery); err != nil {
		return err
	}

	return nil
}
