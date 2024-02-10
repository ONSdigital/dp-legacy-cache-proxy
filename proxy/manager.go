package proxy

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-legacy-cache-proxy/response"
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
	babbageResponse, err := client.Do(proxyReq)

	if err != nil {
		log.Error(ctx, "error sending the proxy request", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer func() {
		if closeErr := babbageResponse.Body.Close(); closeErr != nil {
			log.Error(ctx, "error closing the response body", closeErr)
		}
	}()

	response.WriteResponse(ctx, w, babbageResponse, req)
}
