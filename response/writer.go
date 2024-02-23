package response

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/log.go/v2/log"
)

func WriteResponse(ctx context.Context, w http.ResponseWriter, babbageResponse *http.Response, req *http.Request, cfg *config.Config) {
	if !isGetOrHead(req.Method) {
		writeUnmodifiedResponse(ctx, w, babbageResponse)
	} else if !isCacheableStatusCode(babbageResponse.StatusCode) {
		writeUnmodifiedResponse(ctx, w, babbageResponse)
	} else if cacheControl := babbageResponse.Header.Get("Cache-Control"); !shouldCalculateMaxAge(cacheControl) {
		writeUnmodifiedResponse(ctx, w, babbageResponse)
	} else {
		maxAgeInSeconds := maxAge(ctx, req.RequestURI, cfg)
		writeResponseWithMaxAge(ctx, w, babbageResponse, maxAgeInSeconds)
	}
}

func writeResponse(ctx context.Context, w http.ResponseWriter, babbageResponse *http.Response, overrideHeaders map[string]string) {
	// Copy the Babbage response's headers
	for name, values := range babbageResponse.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set any new headers or overwrite existing
	for name, value := range overrideHeaders {
		w.Header().Set(name, value)
	}

	// Copy the Babbage response's status code
	w.WriteHeader(babbageResponse.StatusCode)

	// Copy the Babbage response's body
	if _, err := io.Copy(w, babbageResponse.Body); err != nil {
		log.Error(ctx, "error copying the proxy response's body", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func writeUnmodifiedResponse(ctx context.Context, w http.ResponseWriter, babbageResponse *http.Response) {
	noAdditionalHeaders := map[string]string{}
	writeResponse(ctx, w, babbageResponse, noAdditionalHeaders)
}

func writeResponseWithMaxAge(ctx context.Context, w http.ResponseWriter, babbageResponse *http.Response, maxAge int) {
	overrideHeaders := make(map[string]string)

	// Get the original Cache-Control value and modify it to include max-age
	originalCacheControl := babbageResponse.Header.Get("Cache-Control")
	if originalCacheControl != "" {
		overrideHeaders["Cache-Control"] = fmt.Sprintf("%s, max-age=%d", originalCacheControl, maxAge)
	} else {
		overrideHeaders["Cache-Control"] = fmt.Sprintf("max-age=%d", maxAge)
	}

	writeResponse(ctx, w, babbageResponse, overrideHeaders)
}

func isGetOrHead(method string) bool {
	return method == http.MethodGet || method == http.MethodHead
}

func isCacheableStatusCode(statusCode int) bool {
	return statusCode < 300 || statusCode == 404
}

func shouldCalculateMaxAge(cacheControlValue string) bool {
	switch cacheControlValue {
	case "", "private", "public":
		return true
	default:
		return false
	}
}
