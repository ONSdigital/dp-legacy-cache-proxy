package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-legacy-cache-proxy
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	OTBatchTimeout             time.Duration `encconfig:"OTEL_BATCH_TIMEOUT"`
	OTExporterOTLPEndpoint     string        `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTServiceName              string        `envconfig:"OTEL_SERVICE_NAME"`
	BabbageURL                 string        `envconfig:"BABBAGE_URL"`
	LegacyCacheAPIURL          string        `envconfig:"LEGACY_CACHE_API_URL"`
	CacheTimeDefault           time.Duration `envconfig:"CACHE_TIME_DEFAULT"`
	CacheTimeErrored           time.Duration `envconfig:"CACHE_TIME_ERRORED"`
	CacheTimeLong              time.Duration `envconfig:"CACHE_TIME_LONG"`
	CacheTimeShort             time.Duration `envconfig:"CACHE_TIME_SHORT"`
	EnablePublishExpiryOffset  bool          `envconfig:"ENABLE_PUBLISH_EXPIRY_OFFSET"`
	PublishExpiryOffset        time.Duration `envconfig:"PUBLISH_EXPIRY_OFFSET"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   ":29200",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		OTBatchTimeout:             5 * time.Second,
		OTExporterOTLPEndpoint:     "localhost:4317",
		OTServiceName:              "dp-legacy-cache-proxy",
		BabbageURL:                 "http://localhost:8080",
		LegacyCacheAPIURL:          "http://localhost:29100",
		CacheTimeDefault:           15 * time.Minute,
		CacheTimeErrored:           30 * time.Second,
		CacheTimeLong:              4 * time.Hour,
		CacheTimeShort:             10 * time.Second,
		EnablePublishExpiryOffset:  false,
		PublishExpiryOffset:        3 * time.Minute,
	}

	return cfg, envconfig.Process("", cfg)
}
