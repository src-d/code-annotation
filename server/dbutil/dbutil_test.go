package dbutil

import (
	"testing"

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

func TestDBUtil(t *testing.T) {
	suite.Run(t, new(DBUtilSuite))
}
