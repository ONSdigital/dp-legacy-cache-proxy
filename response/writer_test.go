package response

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldCalculateMaxAge(t *testing.T) {
	Convey("Given a series of possible values for the Cache-Control header", t, func() {
		testCases := []struct {
			cacheControl string
			expected     bool
		}{
			{cacheControl: "", expected: true},
			{cacheControl: "private", expected: true},
			{cacheControl: "public", expected: true},
			{cacheControl: "max-age=604800", expected: false},
			{cacheControl: "max-age=604800, must-revalidate", expected: false},
			{cacheControl: "s-maxage=604800", expected: false},
			{cacheControl: "no-cache", expected: false},
			{cacheControl: "no-store", expected: false},
			{cacheControl: "public, max-age=604800", expected: false},
			{cacheControl: "public, max-age=604800, immutable", expected: false},
			{cacheControl: "max-age=604800, stale-while-revalidate=86400", expected: false},
			{cacheControl: "max-age=604800, stale-if-error=86400", expected: false},
			{cacheControl: "must-understand, no-store", expected: false},
			{cacheControl: "no-transform", expected: false},
		}
		Convey("When the 'shouldCalculateMaxAge' function is called", func() {
			for _, tc := range testCases {
				result := shouldCalculateMaxAge(tc.cacheControl)
				Convey(fmt.Sprintf(`Then the result should be "%t" when the Cache-Control header is %q`, tc.expected, tc.cacheControl), func() {
					So(result, ShouldEqual, tc.expected)
				})
			}
		})
	})
}
