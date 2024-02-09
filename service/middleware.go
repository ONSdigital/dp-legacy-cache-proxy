package service

import "net/http"

type NoOpRequestMiddleware struct{}

func (rm NoOpRequestMiddleware) GetMiddlewareFunction() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

type HTTPTestRequestMiddleware struct{}

func (rm HTTPTestRequestMiddleware) GetMiddlewareFunction() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The APIFeature in dp-component-test appends "http://foo" to the request. In production, the scheme and
			// host are not set. This middleware removes them, so that the request looks like it would in production.
			r.URL.Scheme = ""
			r.URL.Host = ""

			next.ServeHTTP(w, r)
		})
	}
}
