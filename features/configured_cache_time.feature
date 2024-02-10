Feature: Configured cache time

  The proxy may alter the Cache-Control header in the Babbage response in order to set the "max-age" value to one of
  four preconfigured values: short, long, errored or default cache time.

  Scenario Outline: Versioned path should have a long cache time
    Given Babbage will send the following response:
      """
      Mock response from Babbage
      """
    When the Proxy receives a GET request for "<versioned-path>"
    Then the response header "Cache-Control" should be "max-age=42"
  Examples:
    | versioned-path                                                                                                                                                                          |
    | /economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1                                                                                              |
    | /chartimage?uri=economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1/30d7d6c2                                                                      |
    | /economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/previous/v1/30d7d6c2/data                                                                                |
    | /file?uri=/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpilemindicesandpricequotes/pricequotesseptember2023/previous/v1/pricequotes202309.xlsx |
    | /file?uri=/economy/inflationandpriceindices/datasets/consumerpriceindices/current/previous/v103/mm23.csv                                                                                |
