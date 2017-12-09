Feature: Chat events

  Scenario: Puzzle notification
    Given that the team is enabled
     When a user is notified of the puzzle PUSSGURKA
     Then a chat message containing "PUSSGURKA" is sent to the same user

  Scenario: Correct solution
    Given that the team is enabled
     When a user solved the puzzle
     Then a chat message containing "Korrekt" is sent to the same user

  Scenario: Incorrect solution
    Given that the team is enabled
     When a user attempts an incorrect word
     Then a chat message containing "Inkorrekt" is sent to the same user