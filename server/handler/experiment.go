package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/src-d/code-annotation/server/service"
)

type experiment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Experiment(db *sql.DB) http.HandlerFunc {
	return render(func(r *http.Request) response {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			return respErr(http.StatusBadRequest, err.Error())
		}

		var exp experiment
		err = db.QueryRow("SELECT id, name FROM experiments WHERE id=?", id).
			Scan(&exp.ID, &exp.Name)
		if err != nil {
			return respErr(http.StatusInternalServerError, err.Error())
		}
		return respOK(exp)
	})
}

type nullString struct {
	sql.NullString
}

func (v *nullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}
	return json.Marshal(nil)
}

type assignment struct {
	ID     int        `json:"id"`
	PairID int        `json:"pairId"`
	Answer nullString `json:"answer"`
}

func Assignments(db *sql.DB) http.HandlerFunc {
	return render(func(r *http.Request) response {
		uID := service.GetUserID(r.Context())
		if uID == 0 {
			return respErr(http.StatusInternalServerError, "no user id in context")
		}

		var count int
		rows, err := db.Query("SELECT COUNT(*) as count FROM assignments WHERE user_id=?", uID)
		if err != nil {
			return respErr(http.StatusInternalServerError, err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				return respErr(http.StatusInternalServerError, err.Error())
			}
		}

		if count == 0 {
			tx, err := db.Begin()
			rows, err = tx.Query("SELECT id FROM file_pairs")
			if err != nil {
				tx.Rollback()
				return respErr(http.StatusInternalServerError, err.Error())
			}
			defer rows.Close()
			for rows.Next() {
				var fpID int
				if err := rows.Scan(&fpID); err != nil {
					tx.Rollback()
					return respErr(http.StatusInternalServerError, err.Error())
				}
				_, err := tx.Exec("INSERT INTO assignments (pair_id, user_id) VALUES (?, ?)", fpID, uID)
				if err != nil {
					tx.Rollback()
					return respErr(http.StatusInternalServerError, err.Error())
				}
			}
			if err := rows.Err(); err != nil {
				tx.Rollback()
				return respErr(http.StatusInternalServerError, err.Error())
			}
			tx.Commit()
		}

		var as []assignment
		rows, err = db.Query("SELECT id, pair_id, answer FROM assignments WHERE user_id=?", uID)
		if err != nil {
			return respErr(http.StatusInternalServerError, err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var a assignment
			if err := rows.Scan(&a.ID, &a.PairID, &a.Answer); err != nil {
				return respErr(http.StatusInternalServerError, err.Error())
			}
			as = append(as, a)
		}
		if err := rows.Err(); err != nil {
			return respErr(http.StatusInternalServerError, err.Error())
		}
		return respOK(as)
	})
}

type assignmentInput struct {
	Answer   string `json:"answer"`
	Duration int    `json:"duration"`
}

func UpdateAssignment(db *sql.DB) http.HandlerFunc {
	return render(func(r *http.Request) response {
		uID := service.GetUserID(r.Context())
		if uID == 0 {
			return respErr(http.StatusInternalServerError, "no user id in context")
		}
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			return respErr(http.StatusBadRequest, err.Error())
		}
		var input assignmentInput
		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			return respErr(http.StatusBadRequest, err.Error())
		}
		_, err = db.Exec("UPDATE assignments SET answer=?, duration=? WHERE id=?", input.Answer, input.Duration, id)
		if err != nil {
			return respErr(http.StatusInternalServerError, err.Error())
		}
		return respOK("ok")
	})
}

type filePair struct {
	ID   int    `json:"id"`
	Diff string `json:"diff"`
}

func FilePair(db *sql.DB) http.HandlerFunc {
	return render(func(r *http.Request) response {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			return respErr(http.StatusBadRequest, err.Error())
		}

		var fp filePair
		err = db.QueryRow("SELECT id, diff FROM file_pairs WHERE id=?", id).
			Scan(&fp.ID, &fp.Diff)
		if err != nil {
			return respErr(http.StatusInternalServerError, err.Error())
		}
		return respOK(fp)
	})
}
