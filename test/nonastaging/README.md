# Nonastaging
Nonastaging is a WebSocket based interface to the Nona system. Its purpose is to simplify
acceptance testing of the Nona system deployed in a staging environment.

It is designed to be a testable replacement of the Slack integration `nonaslack`.

# WebSocket interface
Interact with Nona by sending JSON messages over the WebSocket connection to Nonastaging. Below are
examples for each supported message

Request a puzzle for a given user:

    {"type": "user_requests_puzzle", "user_id": "erik"}

## Receiving messages
When Nona sends a response, a chat message is sent from Nonastaging over the WebSocket connection.
It looks like

    {"type": "chat", "user": "erike", "team": "staging", "text": "This is a chat message."}

