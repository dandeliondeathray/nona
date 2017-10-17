Feature: Notify the user when a new puzzle is available

  Scenario: Notify the user of the current puzzle on request
    Given that a users current puzzle is PUSSGURKA
     When a request is sent for the current puzzle
     Then a puzzle notification is sent for puzzle PUSSGURKA
