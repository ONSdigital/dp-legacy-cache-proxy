package steps

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

type ReleaseCalendarFeature struct {
	Server     *httptest.Server
	Body       string
	StatusCode int
	Headers    map[string]string
}

func NewReleaseCalendarFeature() *ReleaseCalendarFeature {
	f := ReleaseCalendarFeature{
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

func (f *ReleaseCalendarFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^Release calendar will send the following response with status "([^"]*)":$`, f.releaseCalendarWillSendTheFollowingResponseWithStatus)
}

func (f *ReleaseCalendarFeature) releaseCalendarWillSendTheFollowingResponseWithStatus(statusCodeStr string, releaseCalendarBody *godog.DocString) error {
	f.Body = strings.TrimSpace(releaseCalendarBody.Content)

	return f.releaseCalendarWillSetTheHTTPStatusCodeTo(statusCodeStr)
}

func (f *ReleaseCalendarFeature) releaseCalendarWillSetTheHTTPStatusCodeTo(statusCodeStr string) error {
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return err
	}

	f.StatusCode = statusCode

	return nil
}
