package service_test

import (
	"io/ioutil"
	"testing"

	"github.com/src-d/code-annotation/server/service"
	"github.com/stretchr/testify/suite"
)

type DiffSuite struct {
	suite.Suite
}

func (suite *DiffSuite) TestDiff() {
	assert := suite.Assert()
	diff := service.NewDiff()

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

	abDiff, err := diff.Generate("a.txt", "b.txt", a, b)

	assert.NoError(err)
	assert.Equal(ab, abDiff)

	acDiff, err := diff.Generate("a.txt", "c.txt", a, c)

	assert.NoError(err)
	assert.Equal(ac, acDiff)
}

func readFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func TestDiff(t *testing.T) {
	suite.Run(t, new(DiffSuite))
}
