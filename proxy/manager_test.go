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

func TestProxyHandleRequestOK(t *testing.T) {
	Convey("Given a Proxy and a Babbage server", t, func() {
		mockBabbageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("mock-header", "test")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("Mock Babbage Response"))
			if err != nil {
				panic(err)
			}
		}))
		defer mockBabbageServer.Close()

		ctx := context.Background()
		router := mux.NewRouter()
		cfg := &config.Config{BabbageURL: mockBabbageServer.URL}

		legacyCacheProxy := proxy.Setup(ctx, router, cfg)

		Convey("When a request is sent", func() {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test-endpoint", http.NoBody)
			legacyCacheProxy.Router.ServeHTTP(w, r)

			Convey("Then the proxy response should match the Babbage response", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				So(w.Body.String(), ShouldEqual, "Mock Babbage Response")
				So(w.Header().Get("mock-header"), ShouldEqual, "test")
			})
		})
	})
}

func TestProxyHandleRequestError(t *testing.T) {
	Convey("Given a Proxy with an invalid Babbage URL configuration", t, func() {
		ctx := context.Background()
		router := mux.NewRouter()
		cfg := &config.Config{BabbageURL: "invalid-babbage-url"}
		legacyCacheProxy := proxy.Setup(ctx, router, cfg)

		Convey("When a request is sent", func() {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test-endpoint", http.NoBody)
			legacyCacheProxy.Router.ServeHTTP(w, r)

			Convey("Then the proxy should return a 500 Internal Server Error", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
