Feature: A pseudo-randomly generated chain of puzzles

  Scenario: Get a puzzle
    Given a new round
     When a request is made for a puzzle at index 0
     Then a puzzle is returned

  Scenario: Next puzzle is different
    Given a new round
     When a request is made for a puzzle at index 1
     Then it is a different puzzle than the one before
  
  Scenario: No round
    Given no new round for a team
     When a request is made for a puzzle at index 0
     Then no puzzle was found