package proxy

import (
	"context"

	"github.com/gorilla/mux"
)

// Proxy provides a struct to wrap the proxy around
type Proxy struct {
	Router *mux.Router
}

// Setup function sets up the proxy and returns a Proxy
func Setup(ctx context.Context, r *mux.Router) *Proxy {
	proxy := &Proxy{
		Router: r,
	}

	// TODO: remove hello world example handler route
	r.HandleFunc("/hello", HelloHandler(ctx)).Methods("GET")
	return proxy
}
