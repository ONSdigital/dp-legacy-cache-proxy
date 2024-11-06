package response

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"github.com/ONSdigital/log.go/v2/log"
)

var (
	resourceEndpoints = []string{"/chartconfig", "/chartimage", "/embed", "/chart", "/resource", "/generator", "/file", "/export"}

	visualisationsEndpointRegexp          = regexp.MustCompile(`(^/visualisations/[^/]+)/`)
	bulletinsOrArticlesPathRegexp         = regexp.MustCompile(`(.+/(bulletins|articles)(?:/[^/]+){2})`)
	relatedDataRegexp                     = regexp.MustCompile(`/related[Dd]ata$`)
	methodologiesOrQMIsOrAdHocsPathRegexp = regexp.MustCompile(`.+/(methodologies|qmis|adhocs)/([^/]+)`)
	fileNameWithExtensionRegexp           = regexp.MustCompile(`(.*)/[^/]+\.\w+$`)
	timeSeriesPathRegexp                  = regexp.MustCompile(`(.+/timeseries(?:/[^/]+){0,2})`)
	datasetsPathRegexp                    = regexp.MustCompile(`(.+/datasets(?:/[^/]+){0,2})`)
)

func getPagePath(ctx context.Context, uri string) (string, error) {
	log.Info(ctx, "Calculating page path", log.Data{"uri": uri})

	uri = strings.TrimSuffix(uri, "/")

	if pagePath, isVisualisationsEndpoint := resolveVisualisationsEndpoint(uri); isVisualisationsEndpoint {
		return pagePath, nil
	}

	if isResourceEndpoint(uri) {
		uriFromQueryString, err := extractAndDecodeURIFromQueryString(ctx, uri)
		if err != nil {
			return "", err
		}

		uri = uriFromQueryString
	}

	if pagePath, isBulletinsOrArticlesPath := resolveBulletinsOrArticlesPath(uri); isBulletinsOrArticlesPath {
		return pagePath, nil
	}

	if pagePath, isMethodologiesQMIsOrAdHocsPath := resolveMethodologiesQMIsOrAdHocsPath(uri); isMethodologiesQMIsOrAdHocsPath {
		return pagePath, nil
	}

	uri = trimUnneededSuffix(uri)

	if pagePath, isTimeSeriesPath := resolveTimeSeriesPath(uri); isTimeSeriesPath {
		return removeLineChartConfig(pagePath), nil
	}

	if pagePath, isDatasetsPath := resolveDatasetPath(uri); isDatasetsPath {
		return pagePath, nil
	}

	return uri, nil
}

func resolveVisualisationsEndpoint(uri string) (string, bool) {
	match := visualisationsEndpointRegexp.FindStringSubmatch(uri)

	if len(match) == 2 {
		return match[1], true
	}

	return "", false
}

func isResourceEndpoint(uri string) bool {
	for _, endpoint := range resourceEndpoints {
		if strings.HasPrefix(uri, endpoint+"?") {
			return true
		}
	}

	return false
}

func extractAndDecodeURIFromQueryString(ctx context.Context, fullURI string) (string, error) {
	urlStruct, err := url.Parse(fullURI)
	if err != nil {
		log.Error(ctx, "error parsing URI", err, log.Data{"uri": fullURI})
		return "", err
	}

	if urlStruct.Query().Has("uri") {
		uriQueryParam := urlStruct.Query().Get("uri")

		decodedURI, err := url.QueryUnescape(uriQueryParam)
		if err != nil {
			log.Error(ctx, "unable to decode the 'uri' query parameter", err, log.Data{"uri": fullURI})
			return "", err
		}

		return strings.TrimSuffix(decodedURI, "/"), nil
	}

	return fullURI, nil
}

func resolveBulletinsOrArticlesPath(uri string) (string, bool) {
	match := bulletinsOrArticlesPathRegexp.FindStringSubmatch(uri)

	if len(match) == 3 {
		return resolveBulletinsOrArticlesExtensionPath(match[1]), true
	}

	return "", false
}

func resolveBulletinsOrArticlesExtensionPath(uri string) string {
	if strings.HasSuffix(uri, "/previousreleases") {
		return strings.TrimSuffix(uri, "/previousreleases") + "/latest"
	}

	if relatedDataRegexp.MatchString(uri) {
		return relatedDataRegexp.ReplaceAllString(uri, "")
	}

	return uri
}

func resolveMethodologiesQMIsOrAdHocsPath(uri string) (string, bool) {
	match := methodologiesOrQMIsOrAdHocsPathRegexp.FindStringSubmatch(uri)

	if len(match) == 3 {
		return match[0], true
	}

	return "", false
}

func trimUnneededSuffix(uri string) string {
	if strings.HasSuffix(uri, "/data") {
		return strings.TrimSuffix(uri, "/data")
	}

	return trimFileNameWithExtension(uri)
}

func trimFileNameWithExtension(uri string) string {
	match := fileNameWithExtensionRegexp.FindStringSubmatch(uri)

	if len(match) == 2 {
		return match[1]
	}

	return uri
}

func resolveTimeSeriesPath(uri string) (string, bool) {
	match := timeSeriesPathRegexp.FindStringSubmatch(uri)

	if len(match) == 2 {
		return match[1], true
	}

	return uri, false
}

func removeLineChartConfig(uri string) string {
	return strings.TrimSuffix(uri, "/linechartconfig")
}

func resolveDatasetPath(uri string) (string, bool) {
	match := datasetsPathRegexp.FindStringSubmatch(uri)

	if len(match) == 2 {
		return match[1], true
	}

	return uri, false
}
