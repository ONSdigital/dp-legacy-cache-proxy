Feature: Configured cache time

  The proxy may alter the Cache-Control header in the Babbage response in order to set the "max-age" value to one of
  four preconfigured values: short, long, errored or default cache time.

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
