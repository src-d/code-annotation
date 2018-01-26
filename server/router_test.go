package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src-d/code-annotation/server/repository"
	"github.com/src-d/code-annotation/server/service"
	testingTools "github.com/src-d/code-annotation/server/testing-tools"

	"github.com/stretchr/testify/suite"
)

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}

type RouterTestSuite struct {
	testingTools.RouterTestSuite
	router http.Handler
	server *httptest.Server
	oauth  *service.OAuth
}

func (s *RouterTestSuite) SetupSuite() {
	s.oauth = service.NewOAuth("client-id", "client-secret")
	s.router = Router(
		testingTools.NewDummyTestLogger(),
		service.NewJWT("sign-key"),
		s.oauth,
		"oauth-redirect-url",
		&repository.Users{},
		"testdata/build",
	)
}

func (s *RouterTestSuite) SetupTest() {
	s.server = httptest.NewServer(s.router)
}

func (s *RouterTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *RouterTestSuite) GetResponse(method string, path string, body io.Reader) *http.Response {
	response, err := testingTools.GetResponse(method, s.server.URL+path, body)
	if err != nil {
		s.Fail(err.Error())
	}

	return response
}

func (s *RouterTestSuite) TestStatics() {
	indexContent := "index\n"
	whateverContent := "content\n"
	deepContent := "deep content\n"
	staticContent := "static content\n"

	var response *http.Response
	response = s.GetResponse("GET", "/", nil)
	s.AssertResponseBodyStatus(response, 200, indexContent, "index should be served")
	response = s.GetResponse("GET", "/whatever-content.xml", nil)
	s.AssertResponseBodyStatus(response, 200, whateverContent, "whatever-content should be served")
	response = s.GetResponse("GET", "/one-level/deep-content.ext", nil)
	s.AssertResponseBodyStatus(response, 200, deepContent, "deep-content should be served")
	response = s.GetResponse("GET", "/non-folder/does-not-exist.ext", nil)
	s.AssertResponseBodyStatus(response, 200, indexContent, "when content does not exist, index should be served")
	response = s.GetResponse("GET", "/does-not-exist.ext", nil)
	s.AssertResponseBodyStatus(response, 200, indexContent, "when content does not exist, index should be served")
	response = s.GetResponse("GET", "/static/static-thing.css", nil)
	s.AssertResponseBodyStatus(response, 200, staticContent, "existent static content should be delivered as it is")
	response = s.GetResponse("GET", "/static/non-existing-static-thing.css", nil)
	s.AssertResponseStatus(response, http.StatusNotFound, "non existent stuff under static directory should return StatusNotFound")
}

func (s *RouterTestSuite) TestLogin() {
	var response *http.Response
	response = s.GetResponse("GET", "/login", nil)
	s.AssertRedirect(response, http.StatusTemporaryRedirect, s.oauth.MakeAuthURL(), "redirect to GH should be performed")
}

func (s *RouterTestSuite) TestOauthCallback() {
	// TODO: FIXME issue: https://github.com/src-d/code-annotation/issues/21
	// TODO: - wrong r.FormValue("state") -> 500

	// TODO: - if given the r.FormValue("code"), the user can not be retrieved from github -> 500

	// TODO: - error in user creation -> 500

	// TODO: - jwt could not be created -> 500

	// TODO: - success -> temporal redirect to fmt.Sprintf("%s/?token=%s", uiDomain, token)
}

func (s *RouterTestSuite) TestUserMe() {
	uID := service.GetUserID(r.Context())
	if uID == 0 {
		errs = append(errs, &serializer.Error{500, "no user id in context", ""})
	}

	u, err := usersRepo.Get(uID)
	if err != nil {
		errs = append(errs, &serializer.Error{404, "user not found", err.Error()})
	}

	resp = serializer.NewUserResponse(u)
	return
}

func (s *RouterTestSuite) TestExperimentDetails() {

	experimentID, err := strconv.Atoi(chi.URLParam(r, "experimentId"))
	if err != nil {
		errs = append(errs, &serializer.Error{500, "Wrong format in experiment ID sent", err.Error()})
		return
	}

	experiment, err := repository.GetExperimentById(experimentID)
	if err != nil {
		errs = append(errs, &serializer.Error{404, "No experiment found", err.Error()})
		return
	}

	resp = serializer.NewExperimentResponse(experiment)
	return

}

func (s *RouterTestSuite) TestFilePair() {
	pairId, err := strconv.Atoi(chi.URLParam(r, "pairId"))
	if err != nil {
		errs = append(errs, &serializer.Error{500, "Wrong format in file-pair ID sent", err.Error()})
		return
	}

	pairFiles, err := repository.GetFilePairFor(pairId)
	if err != nil {
		errs = append(errs, &serializer.Error{404, "No file-pair found", err.Error()})
		return
	}

	resp = serializer.NewFilePairsResponse(pairFiles)
	return
}

func (s *RouterTestSuite) TestAssignmentsForUserAndExperiment() {
	experimentID, err := strconv.Atoi(chi.URLParam(r, "experimentId"))
	if err != nil {
		errs = append(errs, &serializer.Error{500, "Wrong format in experiment ID sent", err.Error()})
		return
	}

	userID, err := strconv.Atoi("1")
	if err != nil {
		errs = append(errs, &serializer.Error{500, "Wrong format in user ID sent", err.Error()})
		return
	}

	assignments, err := repository.GetAssignmentsFor(userID, experimentID)
	if err == repository.ErrNoAssignmentsInitialized {
		assignments, err = repository.CreateAssignmentsFor(userID, experimentID)
	}
	if err != nil {
		errs = append(errs, &serializer.Error{500, "No assignments retrieved", err.Error()})
		return
	}

	resp = serializer.NewAssignmentsResponse(assignments)
	return
	}
}

func (s *RouterTestSuite) TestSaveAssignment() {
	assignmentId, err := strconv.Atoi(chi.URLParam(r, "assignmentId"))
	if err != nil {
		errs = append(errs, &serializer.Error{500, "Wrong format in assignment ID sent", err.Error()})
		return
	}

	var assignmentRequest serializer.AssignmentRequest
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err == nil {
		err = json.Unmarshal(body, &assignmentRequest)
	}
	if err != nil {
		errs = append(errs, &serializer.Error{500, "Payload could not be read", err.Error()})
		return
	}

	if err := repository.UpdateAssignment(assignmentId, assignmentRequest.Answer, assignmentRequest.Duration); err != nil {
		errs = append(errs, &serializer.Error{500, "Answer could not be saved", err.Error()})
		return
	}

	resp = serializer.NewCountResponse(1)
	return
	}
}


