package config

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	os.Clearenv()
	var err error
	var configuration *Config

	Convey("Given an environment with no environment variables set", t, func() {
		Convey("Then cfg should be nil", func() {
			So(cfg, ShouldBeNil)
		})

		Convey("When the config values are retrieved", func() {
			Convey("Then there should be no error returned, and values are as expected", func() {
				configuration, err = Get() // This Get() is only called once, when inside this function
				So(err, ShouldBeNil)
				So(configuration, ShouldResemble, &Config{
					BindAddr:                   ":29200",
					GracefulShutdownTimeout:    5 * time.Second,
					HealthCheckInterval:        30 * time.Second,
					HealthCheckCriticalTimeout: 90 * time.Second,
					OTBatchTimeout:             5 * time.Second,
					OTExporterOTLPEndpoint:     "localhost:4317",
					OTServiceName:              "dp-legacy-cache-proxy",
					BabbageURL:                 "http://localhost:8080",
					LegacyCacheAPIURL:          "http://localhost:29100",
					RelCalURL:                  "http://localhost:27700",
					EnableReleaseCalendar:      false,
					CacheTimeDefault:           15 * time.Minute,
					CacheTimeErrored:           30 * time.Second,
					CacheTimeLong:              4 * time.Hour,
					CacheTimeShort:             10 * time.Second,
					EnablePublishExpiryOffset:  false,
					PublishExpiryOffset:        3 * time.Minute,
					ReadTimeout:                15 * time.Second,
					WriteTimeout:               30 * time.Second,
				})
			})

			Convey("Then a second call to config should return the same config", func() {
				// This achieves code coverage of the first return in the Get() function.
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, cfg)
			})
		})
	})
}
