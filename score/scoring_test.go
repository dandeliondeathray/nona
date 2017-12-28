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

func TestScoring_TwoPlayers_PlayerWithTheHighestScoreIsFirst(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	seed := int64(42)
	persistence := mock.NewFakePersistence()
	response := mock.NewMockResponse(mockCtrl)
	ranking := newRankingMatcher("U2", "U1")

	// Assert
	response.EXPECT().OnPerPlayerScores(gomock.Any(), ranking)

	// Act
	scoring := score.NewScoring(response, persistence)

	scoring.ProduceScores(seed)
	persistence.ResolveAllPlayers(map[game.Player]game.PlayerState{
		game.Player("U1"): game.PlayerState{PuzzleIndex: 17},
		game.Player("U2"): game.PlayerState{PuzzleIndex: 42}})
}

func TestScoring_ThreePlayers_PlayerWithTheHighestScoreIsFirst(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	seed := int64(42)
	persistence := mock.NewFakePersistence()
	response := mock.NewMockResponse(mockCtrl)
	ranking := newRankingMatcher("U2", "U1", "U3")

	// Assert
	response.EXPECT().OnPerPlayerScores(gomock.Any(), ranking)

	// Act
	scoring := score.NewScoring(response, persistence)

	scoring.ProduceScores(seed)
	persistence.ResolveAllPlayers(map[game.Player]game.PlayerState{
		game.Player("U1"): game.PlayerState{PuzzleIndex: 17},
		game.Player("U2"): game.PlayerState{PuzzleIndex: 42},
		game.Player("U3"): game.PlayerState{PuzzleIndex: 1}})
}

// rankingMatcher matches the player rankings against expectations, but does not
// match the scores.
type rankingMatcher struct {
	ranking []game.Player
}

func (m *rankingMatcher) Matches(x interface{}) bool {
	actualRanking, ok := x.([]game.PerPlayerScore)
	previousScore := float64(0.0)
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

		if i != 0 && actualRanking[i].Score > previousScore {
			return false
		} else if i == 0 {
			previousScore = actualRanking[0].Score
		}
	}
	return true
}

func (m *rankingMatcher) String() string {
	return fmt.Sprintf("player ranking %v", m.ranking)
}

func newRankingMatcher(ranking ...game.Player) *rankingMatcher {
	return &rankingMatcher{ranking}
}
