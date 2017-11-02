# SlackMessaging
SlackMessaging is a service that listens to certain events and produces chat
messages, appropriate for consumption by a Slack interface.

## Messages
Consumes these topics

- `nona_PuzzleNotification`

Produces to `nona_<team>_Chat` for configured teams. For instance, for team
"konsulatet", this will produce messages on `nona_konsulatet_Chat`.
