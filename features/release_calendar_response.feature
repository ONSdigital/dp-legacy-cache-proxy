Feature: Proxy returns response from Release Calendar

    The proxy can be called with Release Calendar URLs (for instance, "/releases/adoption"),
    when this happens we want to ensure that the proxy forwards the URL to the Release Calendar. 

  Scenario: Proxy returns response from Release Calendar with stale-while-validate
    Given Release calendar will send the following response with status "200":
      """
      Mock response from Release Calendar
      """
    And config includes STALE_WHILE_REVALIDATE_SECONDS with a value of "31"
    When the Proxy receives a GET request for "/releases/some-path"
    Then the response header "Cache-Control" should be "public, s-maxage=900, max-age=900, stale-while-revalidate=31"
    And the HTTP status code should be "200"
    And I should receive the following response:
      """
      Mock response from Release Calendar
      """


  Scenario: Proxy returns response from Release Calendar without stale-while-validate
    Given Release calendar will send the following response with status "200":
      """
      Mock response from Release Calendar
      """
    And config includes STALE_WHILE_REVALIDATE_SECONDS with a value of "-1"
    When the Proxy receives a GET request for "/releases/some-path"
    Then the response header "Cache-Control" should be "public, s-maxage=900, max-age=900"
    And the HTTP status code should be "200"
    And I should receive the following response:
      """
      Mock response from Release Calendar
      """

