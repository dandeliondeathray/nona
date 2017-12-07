Feature: A pseudo-randomly generated chain of puzzles

  Scenario: Get a puzzle
    Given a new round
     When a request is made for a puzzle at index 0
     Then a puzzle is returned

  Scenario: Next puzzle is different
    Given a new round
     When a request is made for a puzzle at index 1
     Then it is a different puzzle than the one before
      And the solution is a different word
  
  Scenario: No round
    Given no new round for a team
     When a request is made for a puzzle at index 0
     Then no puzzle was found

  Scenario: Service recovery
    Given a new round
     When a request is made for a puzzle at index 5
      And the service goes down and comes back up again
     Then a request for the puzzle at index 5 gives the same answer as before
