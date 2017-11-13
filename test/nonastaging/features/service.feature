Feature: Nonastaging is a WebSocket interface against the nona system

    @wip
    Scenario: A chat message is sent over the WebSocket
         When a chat message is sent to topic nona_staging_Chat
         Then that chat message is received on the WebSocket

    Scenario: User requests a puzzle
         When a user requests a puzzle
         Then a request is sent to nona_UserRequestsPuzzle
