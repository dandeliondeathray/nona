Feature: Chat events

  Scenario: Puzzle notification
    Given that the team is enabled
     When a user is notified of the puzzle PUSSGURKA
     Then a chat message containing PUSSGURKA is sent to the same user
