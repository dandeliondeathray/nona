package score

import (
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
	s.response.OnPerPlayerScores("Pussel-ranking", []game.PerPlayerScore{
		game.PerPlayerScore{Player: game.Player("U1"), Score: 0.0}})
}
