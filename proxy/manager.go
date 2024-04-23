package proxy

import (
	"context"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/dp-legacy-cache-proxy/response"
	"github.com/ONSdigital/log.go/v2/log"
)

func (proxy *Proxy) manage(ctx context.Context, w http.ResponseWriter, req *http.Request, cfg *config.Config) {
	targetURL := getTargetURL(req.URL.String(), cfg)

	proxyReq, err := http.NewRequestWithContext(ctx, req.Method, targetURL, req.Body)

	if err != nil {
		log.Error(ctx, "error creating the proxy request", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Copy headers from original request to proxy request
	proxyReq.Header = req.Header

	client := &http.Client{
		// nolint:revive // param names give context here.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	serviceResponse, err := client.Do(proxyReq)

	if err != nil {
		log.Error(ctx, "error sending the proxy request", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer func() {
		if closeErr := serviceResponse.Body.Close(); closeErr != nil {
			log.Error(ctx, "error closing the response body", closeErr)
		}
	}()

	response.WriteResponse(ctx, w, serviceResponse, req, cfg)
}

func IsReleaseCalendarURL(url string) bool {
	return strings.HasPrefix(url, "/releases/")
}

func getTargetURL(requestURL string, cfg *config.Config) string {
	if IsReleaseCalendarURL(requestURL) && cfg.EnableReleaseCalendar {
		return cfg.RelCalURL + requestURL
	}
	return cfg.BabbageURL + requestURL
}
