Feature: Modification of Cache-Control header to append max-age when Babbage has already set Cache-Control

    There are certain instances in which the Proxy needs to set the max-age directive, but Babbage has already 
    set the Cache-Control header (for instance, to "public" or "private"). In these scenarios, 
    the Proxy has to set the max-age directive without overwriting the existing headers.

  Scenario Outline: Proxy adjusts Cache-Control header when Babbage sets Cache-Control to "public" or "private"
    Given Babbage will send the following response:
      """
      Mock response from Babbage
      """
    And Babbage will set the "Cache-Control" header to "<Babbage-Cache-Control>"
    When the Proxy receives a GET request for "/some-path"
    Then the response header "Cache-Control" should be "<Babbage-Cache-Control>, max-age=900"

  Examples:
    | Babbage-Cache-Control |
    | public                |
    | private               |
