package proxy

import (
	"context"
	"io"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
)

func (proxy *Proxy) manage(ctx context.Context, w http.ResponseWriter, req *http.Request, babbageURL string) {
	targetURL := babbageURL + req.URL.String()

	proxyReq, err := http.NewRequestWithContext(ctx, req.Method, targetURL, req.Body)

	if err != nil {
		log.Error(ctx, "error creating the proxy request", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Copy headers from original request to proxy request
	proxyReq.Header = req.Header

	client := &http.Client{}
	resp, err := client.Do(proxyReq)

	if err != nil {
		log.Error(ctx, "error sending the proxy request", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error(ctx, "error closing the response body", closeErr)
		}
	}()

	// Copy headers from proxy response to original response
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set status code of original response to status code of proxy response
	w.WriteHeader(resp.StatusCode)

	// Copy body of proxy response to original response
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Error(ctx, "error copying the proxy response's body", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
