package handler_test

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/src-d/code-annotation/server/handler"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/stretchr/testify/assert"
)

func TestUploadFilePairs(t *testing.T) {
	assert := assert.New(t)

	// create import db
	dbPath := filepath.Join(os.TempDir(), "cat_test_upload_file_pairs_import_db.db")
	defer os.Remove(dbPath)

	importDB, err := sql.Open("sqlite3", dbPath)
	defer importDB.Close()
	if err != nil {
		t.Fatalf("can't create import db for test %s", err)
	}
	sqlQuery, err := ioutil.ReadFile("testdata/import_db.sql")
	if err != nil {
		t.Fatalf("can't read sql fixture %s", err)
	}
	if _, err := importDB.Exec(string(sqlQuery)); err != nil {
		t.Fatalf("can't apply sql fixture %s", err)
	}
	importDB.Close()

	// create db & handler
	db := testDB()
	handler := handler.UploadFilePairs(db)

	req, err := newFileUploadRequest("/experiments/1/file-pairs", nil, "input_db", dbPath)
	req = chiRequest(req, map[string]string{"experimentId": "1"})
	if err != nil {
		t.Fatalf("can't create file upload request %s", err)
	}
	res, err := handler(req)
	assert.Nil(err)

	assert.Equal(serializer.NewFilePairsUploadResponse(2, 0), res)
}

func newFileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", writer.FormDataContentType())
	return r, nil
}
