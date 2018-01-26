package tools

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/suite"
)

type RouterTestSuite struct {
	suite.Suite
}

func (s *RouterTestSuite) AssertResponseBody(resp *http.Response, expectedContent string, msg string) {
	respBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)
	if err != nil {
		s.FailNow("the response body should be readable")
		return
	}

	defer resp.Body.Close()
	s.Equal(expectedContent, string(respBody), msg)
}

func (s *RouterTestSuite) AssertResponseStatus(resp *http.Response, expectedStatus int, msg string) {
	s.Equal(expectedStatus, resp.StatusCode, fmt.Sprintf("status should be %d; %s", expectedStatus, msg))
}

func (s *RouterTestSuite) AssertResponseBodyStatus(resp *http.Response, expectedStatus int, expectedContent string, msg string) {
	s.AssertResponseBody(resp, expectedContent, msg)
	s.AssertResponseStatus(resp, expectedStatus, "")
}

func (s *RouterTestSuite) AssertResponseHeaders(resp *http.Response, headers map[string]string) {
	for key, value := range headers {
		s.Equal(value, resp.Header.Get(key), fmt.Sprintf("header '%s' should match", key))
	}
}

func (s *RouterTestSuite) AssertRedirect(resp *http.Response, expectedStatus int, expectedURL string, msg string) {
	s.AssertResponseStatus(resp, expectedStatus, "status should be defined by the redirect")
	s.AssertResponseHeaders(resp, map[string]string{"Location": expectedURL})
}

func GetResponse(method string, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, fmt.Errorf("it should be possible to build a request; %s", err.Error())
	}

	resp, err := newHTTPClient(false).Do(req)
	if err != nil {
		return nil, fmt.Errorf("the server should answer with a response; %s", err.Error())
	}

	return resp, nil
}

func newHTTPClient(followRedirects bool) *http.Client {
	client := &http.Client{}
	if !followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return client
}
