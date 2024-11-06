package steps

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

type SearchControllerFeature struct {
	Server     *httptest.Server
	Body       string
	StatusCode int
	Headers    map[string]string
}

func NewSearchControllerFeature() *SearchControllerFeature {
	f := SearchControllerFeature{
		Headers: make(map[string]string),
	}

	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		for headerName, headerValue := range f.Headers {
			w.Header().Set(headerName, headerValue)
		}

		w.WriteHeader(f.StatusCode)

		if _, err := w.Write([]byte(f.Body)); err != nil {
			panic(err)
		}
	}))

	return &f
}

func (f *SearchControllerFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^Search Controller will send the following response:$`, f.searchControllerWillSendTheFollowingResponse)
	ctx.Step(`^Search Controller will send the following response with status "([^"]*)":$`, f.searchControllerWillSendTheFollowingResponseWithStatus)
}

func (f *SearchControllerFeature) searchControllerWillSendTheFollowingResponse(babbageBody *godog.DocString) error {
	return f.searchControllerWillSendTheFollowingResponseWithStatus("200", babbageBody)
}

func (f *SearchControllerFeature) searchControllerWillSendTheFollowingResponseWithStatus(statusCodeStr string, searchControllerBody *godog.DocString) error {
	f.Body = strings.TrimSpace(searchControllerBody.Content)

	return f.searchControllerWillSetTheHTTPStatusCodeTo(statusCodeStr)
}

func (f *SearchControllerFeature) searchControllerWillSetTheHTTPStatusCodeTo(statusCodeStr string) error {
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return err
	}

	f.StatusCode = statusCode

	return nil
}
