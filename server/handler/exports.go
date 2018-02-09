package handler

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/src-d/code-annotation/server/dbutil"
	"github.com/src-d/code-annotation/server/model"
	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"

	"github.com/go-chi/chi"
	"github.com/kelseyhightower/envconfig"
)

const tmpl = `
<html>
	<body>
		<br/>
		<p>Make sure you are logged in from the main dashboard.</p>
		<br/>
		<button onclick="goJWT('/exports/actions/create')">Create a new SQLite export</button>
		<ul>
			{{range .}}
			<li><a href="javascript:goJWT('/exports/actions/download/{{.}}')">{{.}}</a></li>
			{{end}}
		</ul>
		<script>
			function goJWT(url) {
				window.location = url+"?jwt_token="+window.localStorage.getItem("token");
			}
		</script>
	</body>
</html>`

var dbConn string
var usersRepo *repository.Users

// InitializeExports adds the export routes to the given chi.Mux
func InitializeExports(uiDomain string, repo *repository.Users, jwt *service.JWT, r *chi.Mux) {
	// This is not reusing the conf read in server.go to avoid refactoring
	var conf struct {
		DBConn      string `envconfig:"DB_CONNECTION" default:"sqlite://./internal.db"`
		ExportsPath string `envconfig:"EXPORTS_PATH" default:"exports"`
	}
	envconfig.MustProcess("", &conf)

	dbConn = conf.DBConn
	exportsPath := conf.ExportsPath

	// Create the dir to store export files
	os.MkdirAll(exportsPath, 0775)

	usersRepo = repo

	r.Route("/exports", func(r chi.Router) {
		r.Get("/", exportList(exportsPath, true))

		r.Route("/actions", func(r chi.Router) {
			r.Use(jwt.Middleware)
			r.Use(requesterRoleAuth)

			r.Get("/create", exportCreate(exportsPath, uiDomain))
			r.Get("/download/{filename}", exportDownload(exportsPath))
		})
	})
}

// requesterRoleAuth checks that the user making the request has a Requester role
func requesterRoleAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := service.GetUserID(r.Context())

		user, err := usersRepo.GetByID(userID)
		if err != nil {
			exportsErr(err, w, r)
			return
		}

		if user == nil || user.Role != model.Requester {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// exportsErr sets the error message in the log and the response
func exportsErr(err error, w http.ResponseWriter, r *http.Request) {
	err = fmt.Errorf("exports handler error: %s", err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
	service.NewLogger().Error(err.Error())
}

// exportList handles requests to /list
func exportList(exportsPath string, useIndexFallback bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := ioutil.ReadDir(exportsPath)
		if err != nil {
			exportsErr(err, w, r)
			return
		}

		var fileNames []string

		for _, file := range files {
			fileNames = append(fileNames, file.Name())
		}

		sort.Sort(sort.Reverse(sort.StringSlice(fileNames)))

		t := template.New("export tmpl")
		t, err = t.Parse(tmpl)
		if err != nil {
			exportsErr(err, w, r)
			return
		}

		t.Execute(w, fileNames)
	}
}

// exportCreate handles requests to /create
func exportCreate(exportsPath string, uiDomain string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		originDB, err := dbutil.Open(dbConn, true)
		if err != nil {
			exportsErr(err, w, r)
			return
		}
		defer originDB.Close()

		filepath := fmt.Sprintf("./%s/%s-export.db",
			exportsPath, time.Now().Format(time.RFC3339))

		if _, err := os.Stat(filepath); err == nil {
			exportsErr(errors.New("DB file already exists"), w, r)
			return
		}

		destDB, err := dbutil.OpenSQLite(filepath, false)
		if err != nil {
			exportsErr(err, w, r)
			return
		}
		defer destDB.Close()

		if err = dbutil.Bootstrap(destDB); err != nil {
			exportsErr(err, w, r)
			return
		}

		if err := dbutil.Copy(originDB, destDB, dbutil.Options{}); err != nil {
			exportsErr(err, w, r)
			return
		}

		service.NewLogger().Info("new SQLite file created: " + filepath)

		url := fmt.Sprintf("%s/exports", uiDomain)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// exportDownload handles requests to /file/
func exportDownload(exportsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := chi.URLParam(r, "filename")
		filepath := exportsPath + "/" + filename

		file, err := os.Open(filepath)
		if err != nil {
			exportsErr(err, w, r)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		io.Copy(w, file)
	}
}
