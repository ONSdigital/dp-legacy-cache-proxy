package proxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-legacy-cache-proxy/config"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProxyHandleRequestOK(t *testing.T) {
	Convey("Given a Proxy and a Babbage server", t, func() {
		mockBabbageServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("mock-header", "test")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("Mock Babbage Response"))
			if err != nil {
				panic(err)
			}
		}))
		defer mockBabbageServer.Close()
		ctx := context.Background()
		router := mux.NewRouter()
		cfg := &config.Config{BabbageURL: mockBabbageServer.URL}

		legacyCacheProxy := Setup(ctx, router, cfg)

		Convey("When a request is sent", func() {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test-endpoint", http.NoBody)
			legacyCacheProxy.Router.ServeHTTP(w, r)

			Convey("Then the proxy response should match the Babbage response", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				So(w.Body.String(), ShouldEqual, "Mock Babbage Response")
				So(w.Header().Get("mock-header"), ShouldEqual, "test")
			})
		})
	})
}

func TestProxyHandleRequestError(t *testing.T) {
	Convey("Given a Proxy with an invalid Babbage URL configuration", t, func() {
		ctx := context.Background()
		router := mux.NewRouter()
		cfg := &config.Config{BabbageURL: "invalid-babbage-url"}
		legacyCacheProxy := Setup(ctx, router, cfg)

		Convey("When a request is sent", func() {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test-endpoint", http.NoBody)
			legacyCacheProxy.Router.ServeHTTP(w, r)

			Convey("Then the proxy should return a 500 Internal Server Error", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestIsReleaseCalendarURL(t *testing.T) {
	Convey("Given a list of release calendar URLs", t, func() {
		releaseCalendarURLs := []string{
			"/releases/greenjobscurrentandupcomingworkmarch2022",
			"/releases/post2019westminsterparliamentaryconstituenciesandseneddelectoralregionsdataenglandandwalescensus2021",
			"/releases/mycollectionpage1",
			"/releases/timespentinnature",
			"/releases/constructionstatisticsgreatbritain2022",
		}

		Convey("When the 'IsReleaseCalendarURL' function is called", func() {
			for _, url := range releaseCalendarURLs {
				isReleaseCalendar := IsReleaseCalendarURL(url)

				Convey("Then it should evaluate to true for "+url, func() {
					So(isReleaseCalendar, ShouldBeTrue)
				})
			}
		})
	})
	Convey("Given a list of babbage URLs", t, func() {
		babbageURLs := []string{
			"/visualisations/dvc1945/seasonalflu/index.html",
			"/generator?uri=/economy/economicoutputandproductivity/output/bulletins/economicactivityandsocialchangeintheukrealtimeindicators/15february2024/426e63a0&format=csv",
			"/economy/inflationandpriceindices/bulletins/producerpriceinflation/latest",
			"/economy/grossdomesticproductgdp/timeseries/abmi/pn2/linechartconfig",
			"/businessindustryandtrade/changestobusiness/mergersandacquisitions/datasets/timeseries/15march2024",
		}

		Convey("When the 'IsReleaseCalendarURL' function is called", func() {
			for _, url := range babbageURLs {
				isReleaseCalendar := IsReleaseCalendarURL(url)

				Convey("Then it should evaluate to false for "+url, func() {
					So(isReleaseCalendar, ShouldBeFalse)
				})
			}
		})
	})
}
