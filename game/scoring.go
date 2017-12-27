package game

//go:generate mockgen -destination=../mock/mock_game_scoring.go -package=mock github.com/dandeliondeathray/nona/game Scoring

// Scoring is the game interface towards whatever scoring mechanisms there are.
type Scoring interface {
	ProduceScores(seed int64)
}
