Feature: API basics
  
  Scenario: Responds to healthchecks
    When I send "GET" request to "/health"
    Then the response code will be 200
    And the body will be "200 OK"
  
  Scenario: Can create a new poll
    Given I have set data buffer to "requests/poll-create-body.json"
    And I have a cookie "session" with value "session-1"
    When I send "POST" request to "/polls"
    Then the response code will be 201
    And the body will look like json "responses/poll-create-body.json" excluding "poll_id"
    And the database contains a poll with name "poll-request-1"

  Scenario: Can read an existing poll
    Given I have run the database script "db-scripts/base-data.sql"
    And I have a cookie "session" with value "session-1"
    When I make a "GET" request to "/polls/b517b6d9-a25e-4b02-97a2-31fec1fd8afc"
    Then the response code will be 200
    And the body will look like "responses/simple-majority.json"

  Scenario: Get 404 for missing poll
    Given I have run the database script "db-scripts/base-data.sql"
    And I have a cookie "session" with value "session-1"
    When I make a "GET" request to "/polls/aaaaaaaa-a25e-4b02-97a2-aaaaaaaaaaaa"
    Then the response code will be 404
    And the body will look like "responses/not-found.json"

  Scenario: Can update an existing poll
    Given I have run the database script "db-scripts/base-data.sql"
    And I have a cookie "session" with value "session-1"
    And I have set data buffer to "requests/simple-majority-updated.json"
    When I send "PUT" request to "/polls/b517b6d9-a25e-4b02-97a2-31fec1fd8afc"
    Then the response code will be 200
    And the body will look like "responses/simple-majority-updated.json"
    And the database contains a poll with name "simple-majority-1000-renamed"

  Scenario: User can see their only see their own vote
    Given I have run the database script "db-scripts/base-data.sql"
    And I have a cookie "session" with value "session-1"
    When I send "GET" request to "/polls/b517b6d9-a25e-4b02-97a2-31fec1fd8afc/vote"
    Then the response code will be 200
    And the body will look like "responses/simple-majority-user1-vote.json"

  