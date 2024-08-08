package response

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/log.go/v2/log"
)

const maxAgeErrorMessage = "error calculating the max-age directive"

var versionedURIRegexp = regexp.MustCompile(`/previous/v\d+`)
var searchPageRegexp = regexp.MustCompile(`\/(allmethodologies|releasecalendar|timeseriestool|datalist|publications|staticlist|topicspecificmethodology|relateddata|alladhocs|publishedrequests)$`)

func maxAge(ctx context.Context, uri string, cfg *config.Config) (int, bool) {
	log.Info(ctx, "Calculating max-age", log.Data{"uri": uri})

	if isLegacyAssetURI(uri) || isOnsURI(uri) || isVersionedURI(uri) {
		return int(cfg.CacheTimeLong.Seconds()), false
	}

	if isSearchPageURI(uri) {
		return int(cfg.CacheTimeShort.Seconds()), false
	}

	pagePath, err := getPagePath(ctx, uri)
	if err != nil {
		log.Error(ctx, maxAgeErrorMessage, err)
		return int(cfg.CacheTimeErrored.Seconds()), false
	}

	releaseTime, statusCode, err := getReleaseTime(pagePath, cfg.LegacyCacheAPIURL)
	if err != nil {
		log.Error(ctx, maxAgeErrorMessage, err)
		return int(cfg.CacheTimeErrored.Seconds()), false
	}

	if statusCode == http.StatusNotFound {
		return int(cfg.CacheTimeDefault.Seconds()), false
	}

	if statusCode != http.StatusOK {
		unexpectedStatusCodeError := fmt.Errorf("unexpected Legacy Cache API status code: %d", statusCode)
		log.Error(ctx, maxAgeErrorMessage, unexpectedStatusCodeError)
		return int(cfg.CacheTimeErrored.Seconds()), false
	}

	if releaseTime.IsZero() {
		return int(cfg.CacheTimeDefault.Seconds()), false
	}

	if releaseTime.After(time.Now()) {
		if calculatedCacheTime := time.Until(releaseTime); calculatedCacheTime < cfg.CacheTimeDefault {
			return int(calculatedCacheTime.Seconds()), true
		}

		return int(cfg.CacheTimeDefault.Seconds()), false
	}

	if cfg.EnablePublishExpiryOffset && wasReleasedRecently(releaseTime, cfg.PublishExpiryOffset) {
		return int(cfg.CacheTimeShort.Seconds()), false
	}

	return int(cfg.CacheTimeDefault.Seconds()), false
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

func isSearchPageURI(uri string) bool {
	urlToTest, err := url.Parse(uri)
	if err != nil {
		return false
	}

	return searchPageRegexp.MatchString(urlToTest.Path)
}
