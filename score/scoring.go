package score

import (
	"sort"

	"github.com/dandeliondeathray/nona/game"
)

type Scoring struct {
	response    game.Response
	persistence game.Persistence
}

func NewScoring(response game.Response, persistence game.Persistence) *Scoring {
	return &Scoring{response, persistence}
}

func (s *Scoring) ProduceScores(seed int64) {
	simpleScoreCalculator := simpleScore{s.response}
	s.persistence.ResolveAllPlayerStates(&simpleScoreCalculator)
}

// simpleScore calculates the score based on how many puzzles a player
// has solved.
type simpleScore struct {
	response game.Response
}

func (s *simpleScore) AllPlayerStatesResolved(states map[game.Player]game.PlayerState) {
	scores := make([]game.PerPlayerScore, len(states))

	i := 0
	for player, state := range states {
		scores[i] = game.PerPlayerScore{
			Player: player,
			Score:  float64(state.PuzzleIndex - 2*state.Skipped)}
		i++
	}

	sort.Slice(scores, func(i, j int) bool { return scores[i].Score >= scores[j].Score })

	s.response.OnPerPlayerScores("Pussel-ranking", []game.PerPlayerScore(scores))
}
