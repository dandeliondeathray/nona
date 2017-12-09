PuzzleResolver
==============
PuzzleResolver is the service that maintains the puzzle state for a given user. When a user
requests a new puzzle, it is the PuzzleResolver which keeps track of what the current puzzle is.

## Messages
Consumes topic `nona_UserRequestsPuzzle`.

Produces to topic `nona_PuzzleNotification`.

## Environment variables
- **SCHEMA_PATH**: Path to where all Avro schema files are stored.
- **KAFKA_BROKERS**: Comma separated list of URLs to Kafka brokers.