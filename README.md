# dp-legacy-cache-proxy

Proxy for handling the cache for pages within the legacy CMS

### Getting started

- Run `make debug` to run application on http://localhost:29200
- Run `make help` to see full list of make targets

### Dependencies

- No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default                | Description                                                                                                                          |
| ---------------------------- | ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| BIND_ADDR                    | :29200                 | The host and port to bind to                                                                                                         |
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                     | The graceful shutdown timeout in seconds (`time.Duration` format)                                                                    |
| HEALTHCHECK_INTERVAL         | 30s                    | Time between self-healthchecks (`time.Duration` format)                                                                              |
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                    | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)                   |
| OTEL_BATCH_TIMEOUT           | 5s                     | Time duration after which a batch will be sent regardless of size (`time.Duration` format)                                           |
| OTEL_EXPORTER_OTLP_ENDPOINT  | localhost:4317         | OpenTelemetry Exporter address                                                                                                       |
| OTEL_SERVICE_NAME            | dp-legacy-cache-proxy  | The name of this service in OpenTelemetry                                                                                            |
| BABBAGE_URL                  | http://localhost:8080  | Babbage address, where all the incoming requests are forwarded to                                                                    |
| LEGACY_CACHE_API_URL         | http://localhost:29100 | Legacy Cache API address                                                                                                             |
| RELEASE_CALENDAR_URL         | http://localhost:27700 | Release calendar frontend controller address                                                                                         |
| ENABLE_RELEASE_CALENDAR      | false                  | Flag to enable `/releases/{uri:.*}` URLs to go through dp-frontend-release-calendar instead.                                         |
| CACHE_TIME_DEFAULT           | 15m                    | Default value for the `max-age` directive of the `Cache-Control` header (`time.Duration` format)                                     |
| CACHE_TIME_ERRORED           | 30s                    | Errored value for the `max-age` directive of the `Cache-Control` header (`time.Duration` format)                                     |
| CACHE_TIME_LONG              | 4h                     | Long value for the `max-age` directive of the `Cache-Control` header (`time.Duration` format)                                        |
| CACHE_TIME_SHORT             | 10s                    | Short value for the `max-age` directive of the `Cache-Control` header (`time.Duration` format)                                       |
| ENABLE_PUBLISH_EXPIRY_OFFSET | false                  | Determines if publish expiry offset is used which enables a shorter cache time for recently published content.                       |
| PUBLISH_EXPIRY_OFFSET        | 3m                     | Period of time after a release in which the proxy needs to return a short value for the `max-age` directive (`time.Duration` format) |

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2024, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
