package proxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given a Proxy instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		proxy := Setup(ctx, r)

		// TODO: remove hello world example handler route test case
		Convey("When created the following routes should have been added", func() {
			// Replace the check below with any newly added routes
			So(hasRoute(proxy.Router, "/hello", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, http.NoBody)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
