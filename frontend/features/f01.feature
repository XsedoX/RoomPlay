Feature: F01 Logging - Logging in to the application

  Scenario: User logs in with valid credentials
    Given I am on the login page
    When I click Log In with Google
    Then I should be redirected to the main menu page