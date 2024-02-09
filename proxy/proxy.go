package proxy

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

// Proxy provides a struct to wrap the proxy around
type Proxy struct {
	Router *mux.Router
}

// Setup function sets up the proxy and returns a Proxy
func Setup(ctx context.Context, r *mux.Router, babbageURL string) *Proxy {
	proxy := &Proxy{
		Router: r,
	}

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		proxy.manage(ctx, w, req, babbageURL)
	})
	return proxy
}
