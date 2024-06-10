package response

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/log.go/v2/log"
)

func WriteResponse(ctx context.Context, w http.ResponseWriter, serviceResponse *http.Response, req *http.Request, cfg *config.Config) {
	if !isGetOrHead(req.Method) {
		writeUnmodifiedResponse(ctx, w, serviceResponse)
	} else if !isCacheableStatusCode(serviceResponse.StatusCode) {
		writeUnmodifiedResponse(ctx, w, serviceResponse)
	} else if cacheControl := serviceResponse.Header.Get("Cache-Control"); !shouldCalculateMaxAge(cacheControl) {
		writeUnmodifiedResponse(ctx, w, serviceResponse)
	} else {
		maxAgeInSeconds := maxAge(ctx, req.RequestURI, cfg)
		writeResponseWithMaxAge(ctx, w, serviceResponse, maxAgeInSeconds, cfg)
	}
}

func writeResponse(ctx context.Context, w http.ResponseWriter, serviceResponse *http.Response, overrideHeaders map[string]string) {
	// Copy the service response's headers
	for name, values := range serviceResponse.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set any new headers or overwrite existing
	for name, value := range overrideHeaders {
		w.Header().Set(name, value)
	}

	// Copy the service response's status code
	w.WriteHeader(serviceResponse.StatusCode)

	buf := make([]byte, 128*1024)

	// Copy the service response's body
	if _, err := io.CopyBuffer(w, serviceResponse.Body, buf); err != nil {
		log.Error(ctx, "error copying the proxy response's body", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func writeUnmodifiedResponse(ctx context.Context, w http.ResponseWriter, serviceResponse *http.Response) {
	noAdditionalHeaders := map[string]string{}
	writeResponse(ctx, w, serviceResponse, noAdditionalHeaders)
}

func writeResponseWithMaxAge(ctx context.Context, w http.ResponseWriter, serviceResponse *http.Response, maxAge int, cfg *config.Config) {
	overrideHeaders := make(map[string]string)

	cacheControl := "public"
	// Get the original Cache-Control value and - if non-blank - use that instead of above
	originalCacheControl := serviceResponse.Header.Get("Cache-Control")
	if originalCacheControl != "" {
		cacheControl = originalCacheControl
	}
	staleWhileRevalidateOption := ""
	if cfg.StaleWhileRevalidateSeconds >= 0 {
		staleWhileRevalidateOption = fmt.Sprintf(", stale-while-revalidate=%d", cfg.StaleWhileRevalidateSeconds)
	}
	overrideHeaders["Cache-Control"] = fmt.Sprintf("%s, s-maxage=%d, max-age=%d%s", cacheControl, maxAge, maxAge, staleWhileRevalidateOption)

	writeResponse(ctx, w, serviceResponse, overrideHeaders)
}

func isGetOrHead(method string) bool {
	return method == http.MethodGet || method == http.MethodHead
}

func isCacheableStatusCode(statusCode int) bool {
	return statusCode < 300 || statusCode == 304 || statusCode == 404
}

func shouldCalculateMaxAge(cacheControlValue string) bool {
	switch cacheControlValue {
	case "", "private", "public":
		return true
	default:
		return false
	}
}
