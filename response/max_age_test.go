package response

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	. "github.com/smartystreets/goconvey/convey"
)

type TestCase struct {
	cacheTime int
	uri       string
}

func TestMaxAge(t *testing.T) {
	Convey("Given a list of URIs and some pre-configured cache time values", t, func() {
		ctx := context.Background()
		defaultCacheTime := 100
		longCacheTime := 9999
		cfg := &config.Config{
			CacheTimeDefault: time.Duration(defaultCacheTime) * time.Second,
			CacheTimeLong:    time.Duration(longCacheTime) * time.Second,
		}

		versionedURIs := []TestCase{
			{longCacheTime, "/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1"},
			{longCacheTime, "/chartimage?uri=economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1/30d7d6c2"},
			{longCacheTime, "/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1/30d7d6c2/data"},
			{longCacheTime, "/file?uri=/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpilemindicesandpricequotes/pricequotesseptember2023/previous/v1/pricequotes202309.xlsx"},
			{longCacheTime, "/file?uri=/economy/inflationandpriceindices/datasets/consumerpriceindices/current/previous/v103/mm23.csv"},
		}

		onsURIs := []TestCase{
			{longCacheTime, "/ons/rel/household-income/the-effects-of-taxes-and-benefits-on-household-income/index.html"},
			{longCacheTime, "/ons/rel/integrated-household-survey/integrated-household-survey/index.html"},
		}

		legacyAssetURIs := []TestCase{
			{longCacheTime, "/img/national-statistics.png"},
			{longCacheTime, "/css/main.css"},
			{longCacheTime, "/scss/some-sass-file.scss"},
			{longCacheTime, "/js/app.js"},
			{longCacheTime, "/fonts/open-sans-regular/OpenSans-Regular-webfont.woff2"},
			{longCacheTime, "/favicon.ico"},
		}

		regularURIs := []TestCase{
			{defaultCacheTime, "/this-uri/does-not-fall-into-any-special-category"},
			{defaultCacheTime, "/employmentandlabourmarket"},
		}

		groupedTestCases := [][]TestCase{
			versionedURIs,
			onsURIs,
			legacyAssetURIs,
			regularURIs,
		}

		Convey("When the 'maxAge' function is called", func() {
			for _, testCases := range groupedTestCases {
				for _, tc := range testCases {
					result := maxAge(ctx, tc.uri, cfg)

					Convey(fmt.Sprintf(`Then it should return a cache time of %d seconds for the following URI: %q`, tc.cacheTime, tc.uri), func() {
						So(result, ShouldEqual, tc.cacheTime)
					})
				}
			}
		})
	})
}
