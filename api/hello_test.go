package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var ctx = context.Background()

// TODO: remove hello world handler test
func TestHelloHandler(t *testing.T) {
	Convey("Given a Hello handler ", t, func() {
		helloHandler := HelloHandler(ctx)

		Convey("when a good response is returned", func() {
			req := httptest.NewRequest("GET", "http://localhost:8080/hello", nil)
			resp := httptest.NewRecorder()

			helloHandler.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, http.StatusOK)
			So(resp.Body.String(), ShouldResemble, `{"message":"Hello, World!"}`)
		})
	})
}
