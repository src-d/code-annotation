package dbutil

import (
	"io/ioutil"
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

func readFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (suite *DBUtilSuite) TestDiff() {
	assert := suite.Assert()

	a, err := readFile("./testdata/a.txt")
	assert.NoError(err)
	b, err := readFile("./testdata/b.txt")
	assert.NoError(err)
	c, err := readFile("./testdata/c.txt")
	assert.NoError(err)

	ab, err := readFile("./testdata/ab.diff")
	assert.NoError(err)
	ac, err := readFile("./testdata/ac.diff")
	assert.NoError(err)

	abDiff, err := diff("a.txt", "b.txt", a, b)

	assert.NoError(err)
	assert.Equal(ab, abDiff)

	acDiff, err := diff("a.txt", "c.txt", a, c)

	assert.NoError(err)
	assert.Equal(ac, acDiff)
}

func TestDBUtil(t *testing.T) {
	suite.Run(t, new(DBUtilSuite))
}
