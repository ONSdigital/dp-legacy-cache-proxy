# dp-legacy-cache-proxy
Proxy for handling the cache for pages within the legacy CMS

### Getting started

* Run `make debug` to run application on http://localhost:29200
* Run `make help` to see full list of make targets

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default               | Description                                                                                                        |
|------------------------------|-----------------------|--------------------------------------------------------------------------------------------------------------------|
| BIND_ADDR                    | :29200                | The host and port to bind to                                                                                       |
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                    | The graceful shutdown timeout in seconds (`time.Duration` format)                                                  |
| HEALTHCHECK_INTERVAL         | 30s                   | Time between self-healthchecks (`time.Duration` format)                                                            |
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                   | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) |
| OTEL_BATCH_TIMEOUT           | 5s                    | Time duration after which a batch will be sent regardless of size (`time.Duration` format)                         |
| OTEL_EXPORTER_OTLP_ENDPOINT  | localhost:4317        | OpenTelemetry Exporter address                                                                                     |
| OTEL_SERVICE_NAME            | dp-legacy-cache-proxy | The name of this service in OpenTelemetry                                                                          |

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2024, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
