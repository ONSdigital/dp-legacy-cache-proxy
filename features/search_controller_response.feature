Feature: Proxy returns response from Search Controller

    The proxy can be called with Search Controller URLs (for instance, "/economy/economicoutputandproductivity/productivitymeasures/articles/gdpandthelabourmarket/previousreleases"),
    when this happens we want to ensure that the proxy forwards the URL to the Search Controller. 

  Background:
    Given config includes ENABLE_SEARCH_CONTROLLER with a value of "true"

  Scenario: Proxy returns response from Search Controller with stale-while-revalidate
    Given Search Controller will send the following response with status "200":
      """
      Mock response from Search Controller
      """
    And config includes STALE_WHILE_REVALIDATE_SECONDS with a value of "31"
    When the Proxy receives a GET request for "/economy/economicoutputandproductivity/productivitymeasures/articles/gdpandthelabourmarket/previousreleases"
    Then the response header "Cache-Control" should be "public, s-maxage=900, max-age=900, stale-while-revalidate=31"
    And the HTTP status code should be "200"
    And I should receive the following response:
      """
      Mock response from Search Controller
      """

  Scenario: Proxy returns response from Search Controller without stale-while-revalidate
    Given Search Controller will send the following response with status "200":
      """
      Mock response from Search Controller
      """
    And config includes STALE_WHILE_REVALIDATE_SECONDS with a value of "-1"
    When the Proxy receives a GET request for "/economy/economicoutputandproductivity/productivitymeasures/articles/gdpandthelabourmarket/previousreleases"
    Then the response header "Cache-Control" should be "public, s-maxage=900, max-age=900"
    And the HTTP status code should be "200"
    And I should receive the following response:
      """
      Mock response from Search Controller
      """
