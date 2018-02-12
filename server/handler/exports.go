package handler

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/src-d/code-annotation/server/dbutil"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/src-d/code-annotation/server/service"

	"github.com/go-chi/chi"
)

// Export contains handlers for file export
type Export struct {
	db          *dbutil.DB
	exportsPath string
}

// NewExport creates new Export
func NewExport(db *dbutil.DB, exportsPath string) *Export {
	// Create the dir to store export files
	os.MkdirAll(exportsPath, 0775)

	return &Export{
		db:          db,
		exportsPath: exportsPath,
	}
}

// List returns a *serializer.Response with a list of files for export
func (h *Export) List(r *http.Request) (*serializer.Response, error) {
	files, err := ioutil.ReadDir(h.exportsPath)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	sort.Sort(sort.Reverse(sort.StringSlice(fileNames)))

	return &serializer.Response{
		Status: http.StatusOK,
		Data:   fileNames,
	}, nil
}

// Create creates new export file and returns a *serializer.Response
// with the name of new file
func (h *Export) Create(r *http.Request) (*serializer.Response, error) {
	filepath := fmt.Sprintf("%s/%s-export.db",
		h.exportsPath, time.Now().Format(time.RFC3339))

	if _, err := os.Stat(filepath); err == nil {
		return nil, errors.New("DB file already exists")
	}

	destDB, err := dbutil.OpenSQLite(filepath, false)
	if err != nil {
		return nil, err
	}
	defer destDB.Close()

	if err = dbutil.Bootstrap(destDB); err != nil {
		return nil, err
	}

	if err := dbutil.Copy(*h.db, destDB, dbutil.Options{}); err != nil {
		return nil, err
	}

	service.NewLogger().Info("new SQLite file created: " + filepath)

	return &serializer.Response{
		Status: http.StatusOK,
		Data:   filepath,
	}, nil
}

// Download writes requested file to response
func (h *Export) Download(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	filepath := h.exportsPath + "/" + filename

	file, err := os.Open(filepath)
	if err == os.ErrNotExist {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		err = fmt.Errorf("exports handler error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		service.NewLogger().Error(err.Error())
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	io.Copy(w, file)
}
