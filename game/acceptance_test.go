package game_test

import (
	"testing"

	"github.com/dandeliondeathray/nona/game"
	"github.com/dandeliondeathray/nona/mock"
	"github.com/golang/mock/gomock"
)

var acceptanceDictionary = []string{"PUSSGURKA", "PARAMETER"}
var oracle = newGameOracle(acceptanceDictionary)

func TestGiveMeCommand_ForANewRound_PuzzleIsReturned(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	player := game.Player("U1")
	response := mock.NewMockResponse(mockCtrl)
	response.EXPECT().OnPuzzleNotification(player, gomock.Any())

	nona := game.NewGame(response, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player)
}

func TestPuzzles_SolveFirstPuzzle_NextPuzzleIsDifferent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	differentPuzzles := differentPuzzlesMatcher{}

	player := game.Player("U1")
	response := mock.NewMockResponse(mockCtrl)
	response.EXPECT().OnPuzzleNotification(player, &differentPuzzles)
	response.EXPECT().OnPuzzleNotification(player, &differentPuzzles)

	nona := game.NewGame(response, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player)
	correctWord := game.Word(oracle.FindASolutionFor(*differentPuzzles.puzzle))
	nona.TryWord(player, correctWord)
	nona.GiveMe(player)
}

func TestTwoPlayers_FirstPuzzle_SameForBothPlayers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	firstPuzzleIsTheSame := samePuzzlesMatcher{}

	player1 := game.Player("U1")
	player2 := game.Player("U2")
	response := mock.NewMockResponse(mockCtrl)
	response.EXPECT().OnPuzzleNotification(player1, &firstPuzzleIsTheSame)
	response.EXPECT().OnPuzzleNotification(player2, &firstPuzzleIsTheSame)

	nona := game.NewGame(response, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player1)
	nona.GiveMe(player2)
}

func TestTwoPlayers_FirstPlayerSolvesIt_SecondPlayersPuzzleIsUnchanged(t *testing.T) {
	// Assert
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	player2PuzzleUnchanged := samePuzzlesMatcher{}
	player1Puzzle := puzzleSaver{}

	player1 := game.Player("U1")
	player2 := game.Player("U2")
	response := mock.NewMockResponse(mockCtrl)

	response.EXPECT().OnPuzzleNotification(player1, &player1Puzzle)
	response.EXPECT().OnPuzzleNotification(player2, &player2PuzzleUnchanged)
	response.EXPECT().OnPuzzleNotification(player2, &player2PuzzleUnchanged)

	// Arrange
	nona := game.NewGame(response, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player1)
	nona.GiveMe(player2)
	correctWord := game.Word(oracle.FindASolutionFor(*player1Puzzle.puzzle))
	nona.TryWord(player1, correctWord)

	// Act
	// Here player1 has solved the first puzzle, but player2's puzzle should still be the first one.
	nona.GiveMe(player2)
}
