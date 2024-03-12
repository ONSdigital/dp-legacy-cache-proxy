package proxy

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/gorilla/mux"
)

// Proxy provides a struct to wrap the proxy around
type Proxy struct {
	Router *mux.Router
}

// Setup function sets up the proxy and returns a Proxy
func Setup(_ context.Context, r *mux.Router, cfg *config.Config) *Proxy {
	proxy := &Proxy{
		Router: r,
	}

	r.PathPrefix("/").Name("Proxy Catch-All").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		proxy.manage(req.Context(), w, req, cfg)
	})
	return proxy
}
