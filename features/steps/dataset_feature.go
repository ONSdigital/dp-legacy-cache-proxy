package steps

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

type DatasetControllerFeature struct {
	Server     *httptest.Server
	Body       string
	StatusCode int
	Headers    map[string]string
}

func NewDatasetControllerFeature() *DatasetControllerFeature {
	f := DatasetControllerFeature{
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

func (f *DatasetControllerFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^Dataset controller will send the following response with status "([^"]*)":$`, f.datasetControllerWillSendTheFollowingResponseWithStatus)
}

func (f *DatasetControllerFeature) datasetControllerWillSendTheFollowingResponseWithStatus(statusCodeStr string, datasetControllerBody *godog.DocString) error {
	f.Body = strings.TrimSpace(datasetControllerBody.Content)

	return f.datasetControllerWillSetTheHTTPStatusCodeTo(statusCodeStr)
}

func (f *DatasetControllerFeature) datasetControllerWillSetTheHTTPStatusCodeTo(statusCodeStr string) error {
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return err
	}

	f.StatusCode = statusCode

	return nil
}
