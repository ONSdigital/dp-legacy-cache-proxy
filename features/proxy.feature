Feature: Proxy
  Scenario: The response from Babbage is sent back to the client unmodified
    Given Babbage sends the following response:
      """
      Mock response from Babbage
      """
    When the Proxy receives a GET request for "/"
    Then I should receive the following response:
      """
      Mock response from Babbage
      """

  Scenario: The response from Babbage is sent back to the client unmodified, including status code and headers
    Given Babbage sends the following response with status "200":
      """
      Some response
      """
    And Babbage sets the "ETag" header to "abc123"
    And Babbage sets the "Referrer-Policy" header to "origin"
    When the Proxy receives a GET request for "/"
    Then I should receive the following response:
      """
      Some response
      """
    And the HTTP status code should be "200"
    And the response header "ETag" should be "abc123"
    And the response header "Referrer-Policy" should be "origin"
