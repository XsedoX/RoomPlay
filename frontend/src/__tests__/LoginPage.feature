Feature: F01 Login Menu

  Scenario: User sees the apps logo
    Given I am on the login page
    When I look at the login page
    Then I should see the apps logo
