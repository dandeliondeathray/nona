package puzzlestore

import "fmt"

// Teams contains a Puzzles instance for each current round.
// Note that there can only be a single active round for a team at any given time.
type Teams struct {
	rounds     map[string]*Puzzles
	dictionary []string
}

// NewTeams createa new Teams struct.
func NewTeams(dictionary []string) *Teams {
	return &Teams{make(map[string]*Puzzles), dictionary}
}

// NewRound starts a round for a team.
// This must be called before GetPuzzle for a given team, or GetPuzzle will return an error.
func (t *Teams) NewRound(team string, seed int64) {
	t.rounds[team] = NewPuzzles(t.dictionary, seed)
}

// GetPuzzle returns the puzzle for a given index and a given team.
func (t *Teams) GetPuzzle(team string, index int) (string, error) {
	puzzles, ok := t.rounds[team]
	if !ok {
		return "", fmt.Errorf("No puzzles for team '%s'", team)
	}
	return puzzles.Get(index), nil
}
