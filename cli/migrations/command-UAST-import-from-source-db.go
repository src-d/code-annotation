package main

import (
	"fmt"
	"log"

	"github.com/src-d/code-annotation/server/dbutil"
)

type uastImportCmd struct {
	commandDesc
	Args struct {
		InternalDBPath string `description:"filepath to the internal SQLite database"`
		SourceDBPath   string `description:"filepath to the SQLite database containing the UAST to import"`
	} `positional-args:"yes" required:"yes"`
}

var uastImportOpts = uastImportCmd{
	commandDesc: commandDesc{
		name:      "uast-import",
		shortDesc: "Import UASTs",
		longDesc:  "Adds UASTs to the sqlite://InternalDBPath database reading from sqlite://SourceDBPath database",
	},
}

// queries
const (
	filePairsWithoutUastQuery = `SELECT
			blob_id_a, uast_a IS NOT NULL as hasUastA,
			blob_id_b, uast_b IS NOT NULL as hasUastB
		FROM file_pairs
		WHERE hasUastA = 0 or hasUastB = 0`
	uastByBlobIDQuery = `SELECT uast_%s
		FROM files
		WHERE blob_id_%s = CAST($1 AS BLOB)
		LIMIT 1`
	updateByBlobIDQuery = `UPDATE file_pairs
		SET uast_%s = $1
		WHERE blob_id_%s = $2 and uast_%s IS NULL`
	indexAddBlobID = `CREATE INDEX blob_id_%s ON file_pairs (blob_id_%s);`
)

// command messages
const (
	uastImportMsgSuccess = `UASTs added into the internal DB:
		- imported UASTs: %d
		- edited rows: %d`
	uastImportMsgError = `Some rows could not be properly updated:
		- UAST not inserted: %d
		- FilePair read errors: %d`
)

type file struct {
	side   string
	blobID string
}

func (c *uastImportCmd) Execute(args []string) error {
	internalDb, err := dbutil.Open(fmt.Sprintf(sqliteDSN, c.Args.InternalDBPath), false)
	if err != nil {
		log.Fatal(err)
	}

	defer internalDb.Close()

	sourceDb, err := dbutil.Open(fmt.Sprintf(sqliteDSN, c.Args.SourceDBPath), false)
	if err != nil {
		log.Fatal(err)
	}

	defer sourceDb.Close()

	fixSourceDb(sourceDb)

	files, fileReadfailures := getFilesToUpdate(internalDb)
	log.Printf("Found %d blobs without UAST", len(files))

	var rowsEdited, uastFails, uastsImported int64
	for _, file := range files {
		affectedRows, err := importUastForBlobID(internalDb, sourceDb, file.side, file.blobID)
		if err != nil {
			log.Println(err)
			uastFails++
			continue
		}

		rowsEdited += affectedRows
		uastsImported++
	}

	log.Printf(uastImportMsgSuccess, uastsImported, rowsEdited)

	if fileReadfailures+uastFails > 0 {
		log.Fatal(fmt.Sprintf(uastImportMsgError, uastFails, fileReadfailures))
	}

	return nil
}

type files map[string]file

func (f *files) add(blobID string, side string, ignore bool) {
	if ignore {
		return
	}

	if _, ok := (*f)[blobID+"_"+side]; !ok {
		(*f)[blobID+"_"+side] = file{side: side, blobID: blobID}
	}
}

func getFilesToUpdate(internalDb dbutil.DB) (map[string]file, int64) {
	rows, err := internalDb.Query(filePairsWithoutUastQuery)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	filesToImport := files{}
	var failures int64
	for rows.Next() {
		var blobIDA, blobIDB string
		var hasUastA, hasUastB int
		err := rows.Scan(&blobIDA, &hasUastA, &blobIDB, &hasUastB)
		if err != nil {
			log.Printf("Failed to read row from internal DB\nerror: %v\n", err)
			failures++
			continue
		}

		filesToImport.add(blobIDA, leftFile, hasUastA == 1)
		filesToImport.add(blobIDB, rightFile, hasUastB == 1)
	}

	return filesToImport, failures
}

func importUastForBlobID(internalDb dbutil.DB, sourceDb dbutil.DB, side string, blobID string) (int64, error) {
	uast, err := getUastByBlobID(sourceDb, side, blobID)
	if err != nil {
		return 0, fmt.Errorf("uast_%s could not be retrieved for blobID#%s; %s", side, blobID, err)
	}

	return setUastToBlobID(internalDb, side, blobID, uast)
}

func getUastByBlobID(sourceDb dbutil.DB, side string, blobID string) ([]byte, error) {
	var uast []byte
	if err := sourceDb.QueryRow(fmt.Sprintf(uastByBlobIDQuery, side, side), blobID).Scan(&uast); err != nil {
		return nil, err
	}

	return uast, nil
}

func setUastToBlobID(internalDb dbutil.DB, side string, blobID string, uast []byte) (int64, error) {
	res, err := internalDb.Exec(fmt.Sprintf(updateByBlobIDQuery, side, side, side), uast, blobID)
	if err != nil {
		return 0, fmt.Errorf("uast_%s could not be saved for blobID#%s; %s", side, blobID, err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return 0, fmt.Errorf("no uast_%s to be imported for blobID#%s", side, blobID)
	}

	return rows, nil
}

func fixSourceDb(sourceDb dbutil.DB) error {
	if _, err := sourceDb.Exec(fmt.Sprintf(indexAddBlobID, leftFile, leftFile)); err != nil {
		return fmt.Errorf("can not create index over blob_id_%s; %s", leftFile, err)
	}

	if _, err := sourceDb.Exec(fmt.Sprintf(indexAddBlobID, rightFile, rightFile)); err != nil {
		return fmt.Errorf("can not create index over blob_id_%s; %s", rightFile, err)
	}

	return nil
}
