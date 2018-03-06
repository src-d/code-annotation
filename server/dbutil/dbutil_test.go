package dbutil

import (
	"database/sql"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DBUtilSuite struct {
	suite.Suite
}

func (suite *DBUtilSuite) TestMD5() {
	hashes := [][]string{
		{"test string", "6f8db599de986fab7a21625b7916589c"},
		{"Здравствуйте", "66a2e20820c3e976765ccb17b1b7adca"},
		{"a multiline\nstring!", "9beef1614897510967755a19341e730d"}}

	for _, hash := range hashes {
		assert.Equal(suite.T(), hash[1], md5hash(hash[0]))
	}
}

func (suite *DBUtilSuite) TestImportFiles() {
	assert := assert.New(suite.T())

	originDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		suite.T().Fatalf("can't create original db for test %s", err)
	}
	sqlQuery, err := ioutil.ReadFile("testdata/import_db.sql")
	if err != nil {
		suite.T().Fatalf("can't read sql fixture %s", err)
	}
	if _, err := originDB.Exec(string(sqlQuery)); err != nil {
		suite.T().Fatalf("can't apply sql fixture %s", err)
	}

	originDBWrapper := DB{
		DB:     originDB,
		Driver: Sqlite,
	}

	destDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		suite.T().Fatalf("can't create destination db for test %s", err)
	}
	destDBWrapper := DB{
		DB:     destDB,
		Driver: Sqlite,
	}
	if err = Bootstrap(destDBWrapper); err != nil {
		suite.T().Fatalf("can't bootstrap destination db for test %s", err)
	}
	if Initialize(destDBWrapper); err != nil {
		suite.T().Fatalf("can't initialize destination db for test %s", err)
	}

	success, failures, err := ImportFiles(originDBWrapper, destDBWrapper, Options{Logger: logrus.StandardLogger()}, 1)
	assert.NoError(err)
	assert.Equal(int64(2), success)
	assert.Equal(int64(0), failures)
}

func TestDBUtil(t *testing.T) {
	suite.Run(t, new(DBUtilSuite))
}
