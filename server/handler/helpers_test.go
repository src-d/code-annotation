package handler_test

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	"github.com/src-d/code-annotation/server/dbutil"
	"github.com/src-d/code-annotation/server/service"
)

func testDB() *dbutil.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	dbWrapper := dbutil.DB{
		DB:     db,
		Driver: dbutil.Sqlite,
	}
	if err = dbutil.Bootstrap(dbWrapper); err != nil {
		panic(err)
	}
	if dbutil.Initialize(dbWrapper); err != nil {
		panic(err)
	}
	return &dbWrapper
}

func chiRequest(req *http.Request, params map[string]string) *http.Request {
	ctx := lg.WithLoggerContext(req.Context(), logrus.StandardLogger())

	c := chi.NewRouteContext()
	if params != nil {
		for name, value := range params {
			c.URLParams.Add(name, value)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, c)
		}
	}

	return req.WithContext(ctx)
}

func reqWithUser(req *http.Request, userID int) *http.Request {
	ctx := service.SetUserID(req.Context(), userID)
	return req.WithContext(ctx)
}
