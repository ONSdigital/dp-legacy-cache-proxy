# dp-legacy-cache-proxy

The proxy handles the cache for pages within the legacy CMS and sits between the Frontend Router and Babbage. It
receives requests from dp-frontend-router and redirects them to Babbage with the correct Cache-Control header.
All the requests that users make to Babbage (or any other services that rely on the Legacy Cache API for caching
purposes, like the Release Calendar) will go through this Proxy first. When Babbage sends the response back to the user,
the Proxy will intercept it and decide whether it needs to set the max-age directive in the Cache-Control header to an
appropriate value.

## Getting started

- Run `make debug` to run application on [http://localhost:29200](http://localhost:29200)
- Run `make help` to see full list of make targets

### Dependencies

- No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable           | Default                  | Description                                                                                                                                 
|--------------------------------|--------------------------|---------------------------------------------------------------------------------------------------------------------------------------------
| BIND_ADDR                      | :29200                   | The host and port to bind to                                                                                                                
| GRACEFUL_SHUTDOWN_TIMEOUT      | 5s                       | The graceful shutdown timeout in seconds (`time.Duration` format)                                                                           
| HEALTHCHECK_INTERVAL           | 30s                      | Time between self-healthchecks (`time.Duration` format)                                                                                     
| HEALTHCHECK_CRITICAL_TIMEOUT   | 90s                      | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)                          
| HTTP_MAX_CONNECTIONS         | 0                      | Limit the number of concurrent http connections (0 = unlimited)                                                                      |
| OTEL_BATCH_TIMEOUT             | 5s                       | Time duration after which a batch will be sent regardless of size (`time.Duration` format)                                                  
| OTEL_EXPORTER_OTLP_ENDPOINT    | localhost:4317           | OpenTelemetry Exporter address                                                                                                              
| OTEL_SERVICE_NAME              | dp-legacy-cache-proxy    | The name of this service in OpenTelemetry                                                                                                   
| OTEL_ENABLED                   | false                    | Turn OTEL on / off                                                                                                                          
| BABBAGE_URL                    | `http://localhost:8080`  | Babbage address, where all the incoming requests are forwarded to                                                                           
| LEGACY_CACHE_API_URL           | `http://localhost:29100` | Legacy Cache API address                                                                                                                    
| RELEASE_CALENDAR_URL           | `http://localhost:27700` | Release calendar frontend controller address                                                                                                
| CACHE_TIME_DEFAULT             | 15m                      | Default value for the `max-age` directive of the `Cache-Control` header (`time.Duration` format)                                            
| CACHE_TIME_ERRORED             | 30s                      | Errored value for the `max-age` directive of the `Cache-Control` header (`time.Duration` format)                                            
| CACHE_TIME_LONG                | 4h                       | Long value for the `max-age` directive of the `Cache-Control` header (`time.Duration` format)                                               
| CACHE_TIME_SHORT               | 10s                      | Short value for the `max-age` directive of the `Cache-Control` header (`time.Duration` format)                                              
| ENABLE_PUBLISH_EXPIRY_OFFSET   | false                    | Determines if publish expiry offset is used which enables a shorter cache time for recently published content.                              
| PUBLISH_EXPIRY_OFFSET          | 3m                       | Period of time after a release in which the proxy needs to return a short value for the `max-age` directive (`time.Duration` format)        
| READ_TIMEOUT                   | 15s                      | Maximum time the server will wait for a client to send a complete request.                                                                  
| WRITE_TIMEOUT                  | 30s                      | Maximum time the server will wait while trying to write a response to the client                                                            
| STALE_WHILE_REVALIDATE_SECONDS | -1                       | If non-negative, add the `stale-while-revalidate` option (using this number as the *seconds* value) to any `Cache-control` header responses 

## Auto-Deployment of secrets

Functionality has been added to the nomad plan so that when the secrets are deployed to Vault, this will automatically
cause Nomad to trigger a redeployment of the application to pick up the new secrets. Please note that this functionality
does not appear to work with the current nomad/vault versions, but if these are upgraded it may then become functional.

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2024, Office for National Statistics [https://www.ons.gov.uk](https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
