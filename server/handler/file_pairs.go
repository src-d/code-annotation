package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pressly/lg"
	"github.com/src-d/code-annotation/server/dbutil"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/serializer"
	"github.com/src-d/code-annotation/server/service"
)

// GetFilePairDetails returns a function that returns a *serializer.Response
// with the details of the requested FilePair
func GetFilePairDetails(repo *repository.FilePairs, diff *service.Diff) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		pairID, err := urlParamInt(r, "pairId")
		if err != nil {
			return nil, err
		}

		filePair, err := repo.GetByID(pairID)
		if err != nil {
			return nil, err
		}

		if filePair == nil {
			return nil, serializer.NewHTTPError(http.StatusNotFound, "no file-pair found")
		}

		var preprocessors []service.DiffPreprocessorFunc

		if r.URL.Query().Get("showInvisible") == "1" {
			preprocessors = append(preprocessors, service.ReplaceInvisible)
		}

		diffString, err := diff.Generate(
			filePair.Left.Path,
			filePair.Right.Path,
			filePair.Left.Content,
			filePair.Right.Content,
			preprocessors...,
		)
		if err != nil {
			return nil, err
		}

		leftLOC := len(strings.Split(filePair.Left.Content, "\n"))
		rightLOC := len(strings.Split(filePair.Right.Content, "\n"))

		return serializer.NewFilePairResponse(filePair, diffString, leftLOC, rightLOC), nil
	}
}

// GetFilePairs returns a function that returns a *serializer.Response
// with the list of file pairs for the given experiment ID
func GetFilePairs(repo *repository.FilePairs) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		experimentID, err := urlParamInt(r, "experimentId")
		if err != nil {
			return nil, err
		}

		filePairs, err := repo.GetAll(experimentID)
		if err != nil {
			return nil, err
		}

		return serializer.NewListFilePairsResponse(filePairs), nil
	}
}

// UploadFilePairs returns a function that imports file pair from import db file to the experiment
func UploadFilePairs(db *dbutil.DB) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		experimentID, err := urlParamInt(r, "experimentId")
		if err != nil {
			return nil, err
		}

		file, _, err := r.FormFile("input_db")
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		defer file.Close()

		// need to save on disk to be able to open using sql.Open
		tmpfile, err := ioutil.TempFile("", "input_db")
		if err != nil {
			return nil, fmt.Errorf("can't open tmp file for db %s", err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := io.Copy(tmpfile, file); err != nil {
			return nil, fmt.Errorf("can't copy content to tmp db file %s", err)
		}
		if err := tmpfile.Close(); err != nil {
			return nil, fmt.Errorf("can't close tmp db file %s", err)
		}

		inputDB, err := dbutil.OpenSQLite(tmpfile.Name(), false)
		if err != nil {
			return nil, fmt.Errorf("can't open input db %s", err)
		}

		success, failures, err := dbutil.ImportFiles(
			inputDB, *db, dbutil.Options{Logger: lg.RequestLog(r)}, experimentID)
		if err != nil {
			return nil, err
		}

		return serializer.NewFilePairsUploadResponse(success, failures), nil
	}
}
