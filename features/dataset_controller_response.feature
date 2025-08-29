Feature: Proxy returns response from Dataset Controller

    The proxy can be called with Dataset Controller URLs (for instance, "/economy/economicoutputandproductivity/output/datasets/systemaveragepricesapofgas"),
    when this happens we want to ensure that the proxy forwards the URL to the Dataset Controller.

  Scenario: Proxy returns response from Dataset Controller with stale-while-revalidate
    Given Dataset controller will send the following response with status "200":
      """
      Mock response from Dataset Controller
      """
    And config includes STALE_WHILE_REVALIDATE_SECONDS with a value of "31"
    And I set the "Ons-Page-Type" header to "dataset_landing_page"
    When the Proxy receives a GET request for "/some-taxonomy/datasets/some-dataset"
    Then the response header "Cache-Control" should be "public, s-maxage=900, max-age=900, stale-while-revalidate=31"
    And the HTTP status code should be "200"
    And I should receive the following response:
      """
      Mock response from Dataset Controller
      """


  Scenario: Proxy returns response from Dataset Controller without stale-while-revalidate
    Given Dataset controller will send the following response with status "200":
      """
      Mock response from Dataset Controller
      """
    And config includes STALE_WHILE_REVALIDATE_SECONDS with a value of "-1"
    And I set the "Ons-Page-Type" header to "dataset_landing_page"
    When the Proxy receives a GET request for "/some-taxonomy/datasets/some-dataset"
    Then the response header "Cache-Control" should be "public, s-maxage=900, max-age=900"
    And the HTTP status code should be "200"
    And I should receive the following response:
      """
      Mock response from Dataset Controller
      """

