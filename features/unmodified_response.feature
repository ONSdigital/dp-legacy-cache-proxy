Feature: Unmodified response

  Under certain circumstances, the proxy will return the Babbage response to the client completely unmodified. This
  means that the body, headers and status code will be exactly the same and will remain unchanged.

  Scenario Outline: The request method is not GET or HEAD
    Given Babbage will send the following response with status "200":
      """
      Mock response
      """
    And Babbage will set the "ETag" header to "abc123"
    And Babbage will set the "Referrer-Policy" header to "origin"
    When the Proxy receives a <request-method> request for "/"
      """
      Mock request body
      """
    Then I should receive the same, unmodified response from Babbage
  Examples:
    | request-method |
    | POST           |
    | PUT            |
    | PATCH          |
    | DELETE         |

  Scenario Outline: The status code from Babbage is not cacheable
    Given Babbage will send the following response with status "<status-code>":
      """
      Mock response
      """
    And Babbage will set the "ETag" header to "abc123"
    And Babbage will set the "Referrer-Policy" header to "origin"
    When the Proxy receives a GET request for "/"
    Then I should receive the same, unmodified response from Babbage
  Examples:
    | status-code |
    | 300         |
    | 303         |
    | 305         |
    | 306         |
    | 400         |
    | 401         |
    | 402         |
    | 403         |
    | 405         |
    | 406         |
    | 407         |
    | 408         |
    | 409         |
    | 410         |
    | 411         |
    | 412         |
    | 413         |
    | 414         |
    | 415         |
    | 416         |
    | 417         |
    | 418         |
    | 421         |
    | 422         |
    | 423         |
    | 424         |
    | 425         |
    | 426         |
    | 428         |
    | 429         |
    | 431         |
    | 451         |
    | 500         |
    | 501         |
    | 502         |
    | 503         |
    | 504         |
    | 505         |
    | 506         |
    | 507         |
    | 508         |
    | 510         |
    | 511         |

  Scenario Outline: The response from Babbage is not cacheable (based on its Cache-Control header)
    Given Babbage will send the following response with status "200":
      """
      Mock response
      """
    And Babbage will set the "Cache-Control" header to "<cache-control>"
    And Babbage will set the "X-Some-Header" header to "some-value"
    When the Proxy receives a GET request for "/some-url"
    Then I should receive the same, unmodified response from Babbage
  Examples:
    | cache-control          |
    | max-age=123            |
    | no-cache               |
    | no-store               |
    | public, max-age=112233 |
