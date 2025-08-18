Feature: F01 Login Menu

  Scenario: User sees the apps logo
    Given I am on the login page
    When I click Log In with Google
    Then I should be redirected to the main menu page