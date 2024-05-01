package response

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMaxAgeLongCacheTime(t *testing.T) {
	Convey("Given a list of URIs and a pre-configured long cache time", t, func() {
		ctx := context.Background()
		const longCacheTime = 9999
		cfg := &config.Config{
			CacheTimeLong: time.Duration(longCacheTime) * time.Second,
		}

		versionedURIs := []string{
			"/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1",
			"/chartimage?uri=economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1/30d7d6c2",
			"/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1/30d7d6c2/data",
			"/file?uri=/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpilemindicesandpricequotes/pricequotesseptember2023/previous/v1/pricequotes202309.xlsx",
			"/file?uri=/economy/inflationandpriceindices/datasets/consumerpriceindices/current/previous/v103/mm23.csv",
		}

		onsURIs := []string{
			"/ons/rel/household-income/the-effects-of-taxes-and-benefits-on-household-income/index.html",
			"/ons/rel/integrated-household-survey/integrated-household-survey/index.html",
		}

		legacyAssetURIs := []string{
			"/img/national-statistics.png",
			"/css/main.css",
			"/scss/some-sass-file.scss",
			"/js/app.js",
			"/fonts/open-sans-regular/OpenSans-Regular-webfont.woff2",
			"/favicon.ico",
		}

		groupedTestCases := [][]string{
			versionedURIs,
			onsURIs,
			legacyAssetURIs,
		}

		Convey("When the 'maxAge' function is called", func() {
			for _, testCases := range groupedTestCases {
				for _, uri := range testCases {
					result := maxAge(ctx, uri, cfg)

					Convey("Then it should return a long cache time for the following URI: "+uri, func() {
						So(result, ShouldEqual, longCacheTime)
					})
				}
			}
		})
	})
}

func TestMaxAgeShortCacheTime(t *testing.T) {
	Convey("Given a list of URIs and a pre-configured short cache time", t, func() {
		ctx := context.Background()
		const shortCacheTime = 5
		cfg := &config.Config{
			CacheTimeShort: time.Duration(shortCacheTime) * time.Second,
		}

		searchURIs := []string{
			"/releasecalendar",
			"/releasecalendar?view=upcoming",
			"/publications",
			"/economy/datalist",
			"/business/anotherbusiness/allmethodologies",
			"/timeseriestool",
		}

		Convey("When the 'maxAge' function is called", func() {
			for _, uri := range searchURIs {
				result := maxAge(ctx, uri, cfg)

				Convey("Then it should return a short cache time for the following URI: "+uri, func() {
					So(result, ShouldEqual, shortCacheTime)
				})
			}
		})
	})
}

func TestMaxAgeInteractionWithLegacyCacheAPI(t *testing.T) {
	Convey("Given a Legacy Cache API and some pre-configured cache time values", t, func() {
		ctx := context.Background()

		mux := http.NewServeMux()
		mockLegacyCacheAPI := httptest.NewServer(mux)
		defer mockLegacyCacheAPI.Close()

		setMockResponseBody := func(body string) {
			mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(body))
			})
		}

		setMockResponseStatusCode := func(statusCode int) {
			mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(statusCode)
			})
		}

		setMockResponseWithReleaseTime := func(releaseTime time.Time) {
			mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				body := fmt.Sprintf(`{"_id": "7fadfea5c8372c59c0d20599ff95b42a", "path": "/some-valid-path", "release_time": %q}`, releaseTime.Format(time.RFC3339))
				_, _ = w.Write([]byte(body))
			})
		}

		const defaultCacheTime = 100
		const erroredCacheTime = 50
		const shortCacheTime = 3

		const publishExpiryOffset = 15

		cfg := &config.Config{
			LegacyCacheAPIURL:         mockLegacyCacheAPI.URL,
			CacheTimeDefault:          time.Duration(defaultCacheTime) * time.Second,
			CacheTimeErrored:          time.Duration(erroredCacheTime) * time.Second,
			CacheTimeShort:            time.Duration(shortCacheTime) * time.Second,
			EnablePublishExpiryOffset: true,
			PublishExpiryOffset:       time.Duration(publishExpiryOffset) * time.Second,
		}

		Convey("When the 'maxAge' function is called and there is a problem trying to retrieve a Cache Time resource", func() {
			setMockResponseBody("invalid response")
			result := maxAge(ctx, "/some-valid-url", cfg)

			Convey("Then it should return an errored cache time", func() {
				So(result, ShouldEqual, erroredCacheTime)
			})
		})

		Convey("When the 'maxAge' function is called and there is a problem with the API", func() {
			setMockResponseStatusCode(http.StatusInternalServerError)
			result := maxAge(ctx, "/some-valid-url", cfg)

			Convey("Then it should return an errored cache time", func() {
				So(result, ShouldEqual, erroredCacheTime)
			})
		})

		Convey("When the 'maxAge' function is called and the API does not have the requested Cache Time resource", func() {
			setMockResponseStatusCode(http.StatusNotFound)
			result := maxAge(ctx, "/some-valid-url", cfg)

			Convey("Then it should return a default cache time", func() {
				So(result, ShouldEqual, defaultCacheTime)
			})
		})

		Convey("When the 'maxAge' function is called and the requested Cache Time resource does not have a release time", func() {
			setMockResponseBody(`{"_id": "7fadfea5c8372c59c0d20599ff95b42a", "path": "/some-valid-path"}`)
			result := maxAge(ctx, "/some-valid-url", cfg)

			Convey("Then it should return a default cache time", func() {
				So(result, ShouldEqual, defaultCacheTime)
			})
		})

		Convey("When the 'maxAge' function is called and the requested Cache Time resource has a release time", func() {
			Convey("And the release time is in the future", func() {
				Convey("And the release will happen very soon", func() {
					futureReleaseTime := time.Now().Add(30 * time.Second)
					secondsUntilRelease := time.Until(futureReleaseTime).Seconds()
					So(secondsUntilRelease, ShouldBeLessThan, defaultCacheTime)
					setMockResponseWithReleaseTime(futureReleaseTime)
					result := maxAge(ctx, "/some-valid-url", cfg)

					Convey("Then it should return a calculated cache time", func() {
						// Small error threshold (in seconds) to account for result discrepancies due to using an actual
						// time (not mocked) and the tests possibly running slow
						errorThreshold := 3
						So(result, ShouldAlmostEqual, secondsUntilRelease, errorThreshold)
					})
				})

				Convey("And the release will not happen soon", func() {
					futureReleaseTime := time.Now().Add(999999 * time.Hour)
					secondsUntilRelease := time.Until(futureReleaseTime).Seconds()
					So(secondsUntilRelease, ShouldBeGreaterThan, defaultCacheTime)
					setMockResponseWithReleaseTime(futureReleaseTime)
					result := maxAge(ctx, "/some-valid-url", cfg)

					Convey("Then it should return a default cache time", func() {
						So(result, ShouldEqual, defaultCacheTime)
					})
				})
			})

			Convey("And the release time is in the past", func() {
				Convey("And it was released recently", func() {
					pastReleaseTime := time.Now().Add(-3 * time.Second)
					secondsSinceRelease := time.Since(pastReleaseTime).Seconds()
					So(secondsSinceRelease, ShouldBeLessThan, publishExpiryOffset)
					setMockResponseWithReleaseTime(pastReleaseTime)
					result := maxAge(ctx, "/some-valid-url", cfg)

					Convey("Then it should return a short cache time", func() {
						So(result, ShouldEqual, shortCacheTime)
					})
				})

				Convey("And it was not released recently", func() {
					pastReleaseTime := time.Now().Add(-99999999 * time.Second)
					secondsSinceRelease := time.Since(pastReleaseTime).Seconds()
					So(secondsSinceRelease, ShouldBeGreaterThan, publishExpiryOffset)
					setMockResponseWithReleaseTime(pastReleaseTime)
					result := maxAge(ctx, "/some-valid-url", cfg)

					Convey("Then it should return a default cache time", func() {
						So(result, ShouldEqual, defaultCacheTime)
					})
				})
			})
		})
	})
}

func TestMaxAgeWithPublishExpiryOffset(t *testing.T) {
	Convey("Given a Legacy Cache API with some pre-configured cache time values", t, func() {
		ctx := context.Background()

		mux := http.NewServeMux()
		mockLegacyCacheAPI := httptest.NewServer(mux)
		defer mockLegacyCacheAPI.Close()

		setMockResponseWithReleaseTime := func(releaseTime time.Time) {
			mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				body := fmt.Sprintf(`{"_id": "7fadfea5c8372c59c0d20599ff95b42a", "path": "/some-valid-path", "release_time": %q}`, releaseTime.Format(time.RFC3339))
				_, _ = w.Write([]byte(body))
			})
		}

		const defaultCacheTime = 100
		const shortCacheTime = 3
		const publishExpiryOffset = 15

		cfg := &config.Config{
			LegacyCacheAPIURL:   mockLegacyCacheAPI.URL,
			CacheTimeDefault:    time.Duration(defaultCacheTime) * time.Second,
			CacheTimeShort:      time.Duration(shortCacheTime) * time.Second,
			PublishExpiryOffset: time.Duration(publishExpiryOffset) * time.Second,
		}

		Convey("And the requested resource has been recently published", func() {
			veryRecentReleaseTime := time.Now().Add(-10 * time.Second)
			secondsSinceRelease := time.Since(veryRecentReleaseTime).Seconds()
			So(secondsSinceRelease, ShouldBeLessThan, publishExpiryOffset)
			setMockResponseWithReleaseTime(veryRecentReleaseTime)

			Convey("When the Publish Expiry Offset is toggled ON and the 'maxAge' function is called", func() {
				cfg.EnablePublishExpiryOffset = true
				result := maxAge(ctx, "/some-valid-url", cfg)

				Convey("Then it should return a short cache time", func() {
					So(result, ShouldEqual, shortCacheTime)
				})
			})

			Convey("When the Publish Expiry Offset is toggled OFF and the 'maxAge' function is called", func() {
				cfg.EnablePublishExpiryOffset = false
				result := maxAge(ctx, "/some-valid-url", cfg)

				Convey("Then it should return a default cache time", func() {
					So(result, ShouldEqual, defaultCacheTime)
				})
			})
		})
	})
}
