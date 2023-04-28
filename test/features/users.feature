Feature: User specific flows
  
  Scenario: User login with password
    Given I have run the database script "db-scripts/insert-user-1.sql"
    And I have set data buffer to "requests/user1-login.json"
    When I send "POST" request to "/sessions"
    Then the response code will be 204
    And the response has a cookie "session" that is present
  
  Scenario: New user can be registered
    Given I have run the database script "db-scripts/insert-user-1.sql"
    And I have set data buffer to "requests/user2-register.json"
    When I send "POST" request to "/users"
    Then the response code will be 303
    And the response has a header "Location" that matches "/users/\\d+$"

  Scenario: User details can be retrieved
    Given I have run the database script "db-scripts/insert-user-1.sql"
    When I send "GET" request to "/users/b0c88c13-a070-459a-b93a-785cbe39d113"
    Then the response code will be 200
    And the body will look like json "responses/user1-lookup.json" excluding "user_id"