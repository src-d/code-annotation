package service

import (
	"net/http"
	"testing"

	"github.com/src-d/code-annotation/server/model"
	"github.com/stretchr/testify/suite"
)

// roleTestResult is used instead of models.Role to have the "denied" option
type roleTestResult string

const (
	denied        roleTestResult = "denied"
	worker        roleTestResult = "worker"
	requester     roleTestResult = "requester"
	notApplicable roleTestResult = "N/A"
)

const (
	outsiderGroup  = "NONE"
	accessGroup    = "org:accessGroup"
	requesterGroup = "org:requesterGroup"
)

// testRestrictionChecker implements restrictionChecker interface. The mock
// check is done comparing the user's login name with the group restriction
type testRestrictionChecker struct{}

func (checker *testRestrictionChecker) checkAccess(client *http.Client, restriction, login string) error {
	if login != restriction {
		return ErrNoAccess
	}

	return nil
}

func NewTestOAuth(clientID, clientSecret, restrictAccess, restrictRequesterAccess string) *OAuth {
	return &OAuth{
		config:                  nil,
		store:                   nil,
		restrictAccess:          restrictAccess,
		restrictRequesterAccess: restrictRequesterAccess,
		restrictionChecker:      &testRestrictionChecker{},
	}
}

// getTestResult is a helper that performs oauth.setRole
func getTestResult(oauth *OAuth, user GithubUser) roleTestResult {
	err := oauth.setRole(nil, &user)

	if err != nil {
		return denied
	}

	switch user.Role {
	case model.Requester:
		return requester
	case model.Worker:
		return worker
	default:
		return denied
	}
}

type OAuthSuite struct {
	suite.Suite
}

/*
TestGetUser tests role assignment based on the group membership. The tested
logic is the following table:

+-------------------+-----------------------------+-------------+---------------+------------------+
|                   |                             |             |  User in      |  User in         |
|  RESTRICT_ACCESS  |  RESTRICT_REQUESTER_ACCESS  |  Outsider   |  accessGroup  |  requesterGroup  |
+-------------------+-----------------------------+-------------+---------------+------------------+
|  -                |  -                          |  Requester  |  N/A          |  N/A             |
|  -                |  requesterGroup             |  Worker     |  N/A          |  Requester       |
|  accessGroup      |  -                          |  Denied     |  Requester    |  N/A             |
|  accessGroup      |  requesterGroup             |  Denied     |  Worker       |  Requester       |
+-------------------+-----------------------------+-------------+---------------+------------------+
*/
func (suite *OAuthSuite) TestGetUser() {
	var tests = []struct {
		RestrictAccess          string
		RestrictRequesterAccess string
		outsider                roleTestResult
		accessGroupMember       roleTestResult
		requesterGroupMember    roleTestResult
	}{
		{"", "", requester, notApplicable, notApplicable},
		{"", requesterGroup, worker, notApplicable, requester},
		{accessGroup, "", denied, requester, notApplicable},
		{accessGroup, requesterGroup, denied, worker, requester},
	}

	assert := suite.Assert()

	var (
		outsiderUser       = GithubUser{Login: outsiderGroup}
		accessGroupUser    = GithubUser{Login: accessGroup}
		requesterGroupUser = GithubUser{Login: requesterGroup}
	)

	for i, test := range tests {
		oauth := NewTestOAuth("client-id", "client-secret",
			test.RestrictAccess, test.RestrictRequesterAccess)

		assert.Equal(
			test.outsider, getTestResult(oauth, outsiderUser),
			`Failed for test #%v, outsiderUser`, i)

		if test.accessGroupMember != notApplicable {
			assert.Equal(
				test.accessGroupMember, getTestResult(oauth, accessGroupUser),
				`Failed for test #%v, accessGroupUser`, i)
		}

		if test.requesterGroupMember != notApplicable {
			assert.Equal(
				test.requesterGroupMember, getTestResult(oauth, requesterGroupUser),
				`Failed for test #%v, requesterGroupUser`, i)
		}
	}
}

func TestOAuth(t *testing.T) {
	suite.Run(t, new(OAuthSuite))
}
