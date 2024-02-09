package steps

import (
	"net/http"
	"net/http/httptest"
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

	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	ctx.Step(`^Babbage sends the following response:$`, f.babbageSendsTheFollowingResponse)
	ctx.Step(`^Babbage sends the following response with status "([^"]*)":$`, f.babbageSendsTheFollowingResponseWithStatus)
	ctx.Step(`^Babbage sets the "([^"]*)" header to "([^"]*)"$`, f.babbageSetsTheHeaderTo)
}

func (f *BabbageFeature) babbageSendsTheFollowingResponse(babbageBody *godog.DocString) error {
	return f.babbageSendsTheFollowingResponseWithStatus("200", babbageBody)
}

func (f *BabbageFeature) babbageSendsTheFollowingResponseWithStatus(statusCodeStr string, babbageBody *godog.DocString) error {
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return err
	}

	body := strings.TrimSpace(babbageBody.Content)

	f.StatusCode = statusCode
	f.Body = body

	return nil
}

func (f *BabbageFeature) babbageSetsTheHeaderTo(headerName, headerValue string) error {
	f.Headers[headerName] = headerValue

	return nil
}
