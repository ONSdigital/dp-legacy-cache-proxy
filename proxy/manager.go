package proxy

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/dp-legacy-cache-proxy/response"
	"github.com/ONSdigital/log.go/v2/log"
)

func (proxy *Proxy) manage(ctx context.Context, w http.ResponseWriter, req *http.Request, cfg *config.Config) {
	pageType := req.Header.Get("Ons-Page-Type")
	targetURL := getTargetURL(req.URL.String(), pageType, cfg)

	proxyReq, err := http.NewRequestWithContext(ctx, req.Method, targetURL, req.Body)

	if err != nil {
		log.Error(ctx, "error creating the proxy request", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Copy headers from original request to proxy request
	proxyReq.Header = req.Header
	// Also copy Host (header had been removed from original request)
	proxyReq.Host = req.Host

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

func IsReleaseCalendarURL(requestURLstring string) bool {
	return strings.HasPrefix(requestURLstring, "/releases/")
}

func IsSearchControllerURL(requestURLstring string) bool {
	requestURL, err := url.Parse(requestURLstring)
	if err != nil {
		return false
	}
	return (strings.HasSuffix(requestURL.EscapedPath(), "/previousreleases") || strings.HasSuffix(requestURL.EscapedPath(), "/relatedData") || strings.HasSuffix(requestURL.EscapedPath(), "/relateddata"))
}

// IsDatasetLandingPage checks if the page type is a dataset_landing_page
func IsDatasetLandingPage(pageTypeString string) bool {
	return pageTypeString == "dataset_landing_page"
}

func getTargetURL(requestURL, pageType string, cfg *config.Config) string {
	if IsReleaseCalendarURL(requestURL) {
		return cfg.RelCalURL + requestURL
	} else if IsSearchControllerURL(requestURL) && cfg.EnableSearchController {
		return cfg.SearchControllerURL + requestURL
	} else if IsDatasetLandingPage(pageType) {
		return cfg.DatasetControllerURL + requestURL
	}
	return cfg.BabbageURL + requestURL
}
