package response

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetReleaseTime(t *testing.T) {
	Convey("Given a Legacy Cache API", t, func() {
		var writeMockResponse func(w http.ResponseWriter) error
		mockLegacyCacheAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := writeMockResponse(w)
			if err != nil {
				t.Fatal("error setting the mock server response body", err)
			}
		}))
		defer mockLegacyCacheAPI.Close()

		Convey("When 'getReleaseTime' is called and a Cache Time resource with a Release Time is returned from the API", func() {
			cacheTimeResource := `
			{
				"_id": "7fadfea5c8372c59c0d20599ff95b42a",
				"path": "/some-valid-path",
				"collection_id": 123456,
				"release_time": "2024-01-31T01:23:45.678Z"
			}`
			writeMockResponse = setMockResponse(cacheTimeResource, http.StatusOK)

			releaseTime, statusCode, err := getReleaseTime("/some-valid-path", mockLegacyCacheAPI.URL)

			Convey("Then the result is the Release Time and status code with no errors", func() {
				expectedReleaseTime, _ := time.Parse(time.RFC3339, "2024-01-31T01:23:45.678Z")
				So(releaseTime, ShouldEqual, expectedReleaseTime)
				So(statusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
			})
		})

		Convey("When 'getReleaseTime' is called and a Cache Time resource with no Release Time is returned from the API", func() {
			cacheTimeResource := `
			{
				"_id": "7fadfea5c8372c59c0d20599ff95b42a",
				"path": "/some-valid-path"
			}`
			writeMockResponse = setMockResponse(cacheTimeResource, http.StatusOK)

			releaseTime, statusCode, err := getReleaseTime("/some-valid-path", mockLegacyCacheAPI.URL)

			Convey("Then the result is an empty Release Time and status code with no errors", func() {
				So(releaseTime.IsZero(), ShouldBeTrue)
				So(statusCode, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
			})
		})

		Convey("When 'getReleaseTime' is called and a Cache Time resource is not found in the API", func() {
			writeMockResponse = setMockResponse("", http.StatusNotFound)

			releaseTime, statusCode, err := getReleaseTime("/some-valid-path", mockLegacyCacheAPI.URL)

			Convey("Then the result is an empty Release Time and a Not Found status code with no errors", func() {
				So(releaseTime.IsZero(), ShouldBeTrue)
				So(statusCode, ShouldEqual, http.StatusNotFound)
				So(err, ShouldBeNil)
			})
		})

		Convey("When 'getReleaseTime' is called and an unexpected status code is returned from the API", func() {
			writeMockResponse = setMockResponse("", http.StatusBadGateway)

			releaseTime, statusCode, err := getReleaseTime("/some-valid-path", mockLegacyCacheAPI.URL)

			Convey("Then the result is an empty Release Time and the same status code with no errors", func() {
				So(releaseTime.IsZero(), ShouldBeTrue)
				So(statusCode, ShouldEqual, http.StatusBadGateway)
				So(err, ShouldBeNil)
			})
		})

		Convey("When 'getReleaseTime' is called and there is an error with the API", func() {
			writeMockResponse = nil

			_, _, err := getReleaseTime("/some-valid-path", "invalid-API-URL")

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func setMockResponse(body string, statusCode int) func(http.ResponseWriter) error {
	return func(w http.ResponseWriter) error {
		w.WriteHeader(statusCode)
		_, err := fmt.Fprint(w, body)
		return err
	}
}
