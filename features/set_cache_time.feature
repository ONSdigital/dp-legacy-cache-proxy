Feature: Set cache time

  The proxy may alter the Cache-Control header in the Babbage response in order to set the "max-age" directive to one of
  four preconfigured values: short, long, errored or default cache time. It may also be set to a calculated value if
  it is a page that is about to be released.

  Background:
    Given Babbage will send the following response:
      """
      Mock response from Babbage
      """

  Scenario Outline: Versioned URI should have a long cache time
    When the Proxy receives a GET request for "<versioned-uri>"
    Then the response header "Cache-Control" should be "max-age=14400"
  Examples:
    | versioned-uri                                                                                                                                                                           |
    | /economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1                                                                                              |
    | /chartimage?uri=economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1/30d7d6c2                                                                      |
    | /economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1/30d7d6c2/data                                                                                |
    | /file?uri=/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpilemindicesandpricequotes/pricequotesseptember2023/previous/v1/pricequotes202309.xlsx |
    | /file?uri=/economy/inflationandpriceindices/datasets/consumerpriceindices/current/previous/v103/mm23.csv                                                                                |

  Scenario Outline: ONS URI should have a long cache time
    When the Proxy receives a GET request for "<ons-uri>"
    Then the response header "Cache-Control" should be "max-age=14400"
  Examples:
    | ons-uri                                                                                    |
    | /ons/rel/household-income/the-effects-of-taxes-and-benefits-on-household-income/index.html |
    | /ons/rel/integrated-household-survey/integrated-household-survey/index.html                |

  Scenario Outline: Legacy asset URI should have a long cache time
    When the Proxy receives a GET request for "<legacy-asset-uri>"
    Then the response header "Cache-Control" should be "max-age=14400"
  Examples:
    | legacy-asset-uri                                        |
    | /img/national-statistics.png                            |
    | /css/main.css                                           |
    | /scss/some-sass-file.scss                               |
    | /js/app.js                                              |
    | /fonts/open-sans-regular/OpenSans-Regular-webfont.woff2 |
    | /favicon.ico                                            |

  Scenario: Return the errored cache time when the Legacy Cache API returns an error
    Given the Legacy Cache API has an error
    When the Proxy receives a GET request for "/some-path"
    Then the response header "Cache-Control" should be "max-age=30"

  Scenario: Return the default cache time when the Legacy Cache API does not have the requested page
    Given the Legacy Cache API does not have any data for the "/some-path" page
    When the Proxy receives a GET request for "/some-path"
    Then the response header "Cache-Control" should be "max-age=900"

  Scenario: Return the default cache time when the release time is missing
    Given the "/some-path" page does not have a release time
    When the Proxy receives a GET request for "/some-path"
    Then the response header "Cache-Control" should be "max-age=900"

  Scenario: Return the calculated cache time when the release time is in the near future
    Given the "/some-path" page will have a release in the near future
    When the Proxy receives a GET request for "/some-path"
    Then the max-age directive should be calculated, rather than predefined

  Scenario: Return the default cache time when the release time is in the distant future
    Given the "/some-path" page will have a release in the distant future
    When the Proxy receives a GET request for "/some-path"
    Then the response header "Cache-Control" should be "max-age=900"

  Scenario: Return the default cache time when the page was released long ago
    Given the "/some-path" page was released long ago
    When the Proxy receives a GET request for "/some-path"
    Then the response header "Cache-Control" should be "max-age=900"

  Scenario: Return the short cache time when the page was released recently
    Given the "/some-path" page was released recently
    When the Proxy receives a GET request for "/some-path"
    Then the response header "Cache-Control" should be "max-age=10"

  Scenario: Return the default cache time when the page was released recently and the publish expiry offset is disabled
    Given the "/some-path" page was released recently
    And the Proxy has the publish expiry offset disabled
    When the Proxy receives a GET request for "/some-path"
    Then the response header "Cache-Control" should be "max-age=900"
