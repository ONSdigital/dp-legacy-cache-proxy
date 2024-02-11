package proxy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/dp-legacy-cache-proxy/proxy"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given a Proxy instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		cfg := &config.Config{}
		legacyCacheProxy := proxy.Setup(ctx, r, cfg)

		Convey("When created, all HTTP methods should be accepted", func() {
			So(hasRoute(legacyCacheProxy.Router, "/", http.MethodGet), ShouldBeTrue)
			So(hasRoute(legacyCacheProxy.Router, "/", http.MethodPost), ShouldBeTrue)
			So(hasRoute(legacyCacheProxy.Router, "/", http.MethodPut), ShouldBeTrue)
			So(hasRoute(legacyCacheProxy.Router, "/", http.MethodDelete), ShouldBeTrue)
			So(hasRoute(legacyCacheProxy.Router, "/", http.MethodHead), ShouldBeTrue)
			So(hasRoute(legacyCacheProxy.Router, "/", http.MethodConnect), ShouldBeTrue)
			So(hasRoute(legacyCacheProxy.Router, "/", http.MethodOptions), ShouldBeTrue)
			So(hasRoute(legacyCacheProxy.Router, "/", http.MethodTrace), ShouldBeTrue)
			So(hasRoute(legacyCacheProxy.Router, "/", http.MethodPatch), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, http.NoBody)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
