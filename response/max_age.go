package response

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/log.go/v2/log"
)

var versionedURIRegexp = regexp.MustCompile(`/previous/v\d+`)

func maxAge(ctx context.Context, uri string, cfg *config.Config) int {
	log.Info(ctx, "Calculating max-age for "+uri)

	if isLegacyAssetURI(uri) || isOnsURI(uri) || isVersionedURI(uri) {
		return int(cfg.CacheTimeLong.Seconds())
	}

	pagePath := getPagePath(ctx, uri)

	releaseTime, statusCode, err := getReleaseTime(pagePath, cfg.LegacyCacheAPIURL)
	if err != nil {
		log.Error(ctx, "error calculating the max-age directive", err)
		return int(cfg.CacheTimeErrored.Seconds())
	}

	if statusCode == http.StatusNotFound {
		return int(cfg.CacheTimeDefault.Seconds())
	}

	if statusCode != http.StatusOK {
		unexpectedStatusCodeError := fmt.Errorf("unexpected Legacy Cache API status code: %d", statusCode)
		log.Error(ctx, "error calculating the max-age directive", unexpectedStatusCodeError)
		return int(cfg.CacheTimeErrored.Seconds())
	}

	if releaseTime.IsZero() {
		return int(cfg.CacheTimeDefault.Seconds())
	}

	if releaseTime.After(time.Now()) {
		if calculatedCacheTime := time.Until(releaseTime); calculatedCacheTime < cfg.CacheTimeDefault {
			return int(calculatedCacheTime.Seconds())
		}

		return int(cfg.CacheTimeDefault.Seconds())
	}

	if wasReleasedRecently(releaseTime, cfg.PublishExpiryOffset) {
		return int(cfg.CacheTimeShort.Seconds())
	}

	return int(cfg.CacheTimeDefault.Seconds())
}

func isLegacyAssetURI(uri string) bool {
	legacyAssetFolders := []string{"/img/", "/css/", "/scss/", "/js/", "/fonts/"}

	for _, folder := range legacyAssetFolders {
		if strings.HasPrefix(uri, folder) {
			return true
		}
	}

	return uri == "/favicon.ico"
}

func isOnsURI(uri string) bool {
	return strings.HasPrefix(uri, "/ons/")
}

func isVersionedURI(uri string) bool {
	return versionedURIRegexp.MatchString(uri)
}

func wasReleasedRecently(releaseTime time.Time, offset time.Duration) bool {
	return releaseTime.Add(offset).After(time.Now())
}
