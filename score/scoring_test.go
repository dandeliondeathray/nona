package score_test

import (
	"fmt"
	"testing"

	"github.com/dandeliondeathray/nona/game"
	"github.com/dandeliondeathray/nona/mock"
	"github.com/dandeliondeathray/nona/score"
	"github.com/golang/mock/gomock"
)

func TestScoring_OnePlayer_PerPlayerScoreResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	seed := int64(42)
	persistence := mock.NewFakePersistence()
	response := mock.NewMockResponse(mockCtrl)
	ranking := newRankingMatcher("U1")

	// Assert
	response.EXPECT().OnPerPlayerScores(gomock.Any(), ranking)

	// Act
	scoring := score.NewScoring(response, persistence)

	scoring.ProduceScores(seed)
	persistence.ResolveAllPlayers(map[game.Player]game.PlayerState{
		game.Player("U1"): game.PlayerState{PuzzleIndex: 17}})
}

// rankingMatcher matches the player rankings against expectations, but does not
// match the scores.
type rankingMatcher struct {
	ranking []game.Player
}

func (m *rankingMatcher) Matches(x interface{}) bool {
	actualRanking, ok := x.([]game.PerPlayerScore)
	if !ok {
		return false
	}
	if len(m.ranking) != len(actualRanking) {
		return false
	}
	for i := 0; i < len(m.ranking); i++ {
		if m.ranking[i] != actualRanking[i].Player {
			return false
		}
	}
	return true
}

func (m *rankingMatcher) String() string {
	return fmt.Sprintf("rankingMatcher{%v}", m.ranking)
}

func newRankingMatcher(ranking ...game.Player) *rankingMatcher {
	return &rankingMatcher{ranking}
}
