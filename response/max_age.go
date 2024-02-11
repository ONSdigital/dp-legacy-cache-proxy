package response

import (
	"context"
	"regexp"
	"strings"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/ONSdigital/log.go/v2/log"
)

var versionedURIRegexp = regexp.MustCompile(`/previous/v\d+`)

func maxAge(ctx context.Context, uri string, cfg *config.Config) int {
	log.Info(ctx, "Calculating max-age for "+uri)

	if isLegacyAssetURI(uri) || isOnsURI(uri) || isVersionedURI(uri) {
		return int(cfg.CacheTimeLong.Seconds())
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
	matched := versionedURIRegexp.MatchString(uri)

	return matched
}