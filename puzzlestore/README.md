Puzzlestore
===========
Puzzlestore is the service that maintains the puzzle state. When a user
requests a new puzzle, it is the Puzzlestore which keeps track of what
the current puzzle is.

## Messages
Consumes topic `nona_UserRequestsPuzzle`.

Produces to topic `nona_PuzzleNotification`.
