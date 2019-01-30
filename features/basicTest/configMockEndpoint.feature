Feature: load config file and test configuration object's behaviour
  1. to be able to load and parse the given config file (or through byte[] raw source)
  2. able to load back the configs / endpoints of the given data (file or byte[])
  3. able to retrieve the config contents based on method name (e.g. getAuthors) and the http verb (e.g. GET)

  Scenario: basic test
    Given a config file named "sampleMockInstructions.json"
    When successfully loaded and parsed, 5 methods are available
    Then get back the configured method "getAuthorById" having http verb "GET" should return a method definition
    Then the object returned should consists of a "json" response stating id => 13 would get back an author named "Eliza"

  Scenario: removeAuthorById test
    Given a config file named "sampleMockInstructions.json"
    When successfully loaded and parsed, 5 methods are available
    Then get back the configured method "removeAuthorById" having http verb "DELETE" should return a method definition
    Then the object returned should consists of a "xml" response stating id => 1 would get back an author named "Jacky"