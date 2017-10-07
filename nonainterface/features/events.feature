Feature: User interface against the microservices
  Nonainterface has an EVENT queue, which is a thread safe queue on which all messages directed to Slack are put.
  A CHAT MESSAGE is a message directed towards a specific user.
  A NOTIFICATION is a message directed towards a channel, public or private.

  Scenario: A user requests a puzzle
    Given a team konsulatet
     When a user requests a puzzle
     Then a UserRequestsPuzzle is sent to topic nona_konsulatet_UserRequestsPuzzle

  Scenario: Puzzle response
    Given a team konsulatet
     When there is a puzzle response in nona_konsulatet_Chat
     Then a chat message is in the event queue