package steps

import (
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

type BabbageFeature struct {
	Server     *httptest.Server
	Body       string
	StatusCode int
	Headers    map[string]string
}

func NewBabbageFeature() *BabbageFeature {
	f := BabbageFeature{
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

func (f *BabbageFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^Babbage will send the following response:$`, f.babbageWillSendTheFollowingResponse)
	ctx.Step(`^Babbage will send the following response with status "([^"]*)":$`, f.babbageWillSendTheFollowingResponseWithStatus)
	ctx.Step(`^Babbage will set the "([^"]*)" header to "([^"]*)"$`, f.babbageWillSetTheHeaderTo)
	ctx.Step(`^Babbage will set the HTTP status code to "([^"]*)"$`, f.babbageWillSetTheHTTPStatusCodeTo)
}

func (f *BabbageFeature) babbageWillSendTheFollowingResponse(babbageBody *godog.DocString) error {
	return f.babbageWillSendTheFollowingResponseWithStatus("200", babbageBody)
}

func (f *BabbageFeature) babbageWillSendTheFollowingResponseWithStatus(statusCodeStr string, babbageBody *godog.DocString) error {
	f.Body = strings.TrimSpace(babbageBody.Content)

	return f.babbageWillSetTheHTTPStatusCodeTo(statusCodeStr)
}

func (f *BabbageFeature) babbageWillSetTheHeaderTo(headerName, headerValue string) error {
	canonicalHeaderName := textproto.CanonicalMIMEHeaderKey(headerName)
	f.Headers[canonicalHeaderName] = headerValue

	return nil
}

func (f *BabbageFeature) babbageWillSetTheHTTPStatusCodeTo(statusCodeStr string) error {
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return err
	}

	f.StatusCode = statusCode

	return nil
}
