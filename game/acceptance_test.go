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
	persistence := newFakePersistence()

	// Assert
	response.EXPECT().OnPuzzleNotification(player, gomock.Any())

	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player)
	persistence.playerStateResolved(player)
}

func TestPuzzles_SolveFirstPuzzle_NextPuzzleIsDifferent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	differentPuzzles := differentPuzzlesMatcher{}

	player := game.Player("U1")
	response := mock.NewMockResponse(mockCtrl)
	persistence := newFakePersistence()
	response.EXPECT().OnPuzzleNotification(player, &differentPuzzles)
	response.EXPECT().OnPuzzleNotification(player, &differentPuzzles)
	response.EXPECT().OnCorrectWord(gomock.Any(), gomock.Any()).AnyTimes()

	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player)
	persistence.playerStateResolved(player)
	correctWord := game.Word(oracle.FindASolutionFor(*differentPuzzles.puzzle))
	nona.TryWord(player, correctWord)
	persistence.playerStateResolved(player)
	nona.GiveMe(player)
	persistence.playerStateResolved(player)
}

func TestPuzzles_TryIncorrectSolution_CurrentPuzzleIsUnchanged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	puzzleIsUnchanged := samePuzzlesMatcher{}

	player := game.Player("U1")
	response := mock.NewMockResponse(mockCtrl)
	persistence := newFakePersistence()
	response.EXPECT().OnPuzzleNotification(player, &puzzleIsUnchanged)
	response.EXPECT().OnPuzzleNotification(player, &puzzleIsUnchanged)
	response.EXPECT().OnCorrectWord(gomock.Any(), gomock.Any()).AnyTimes()

	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player)
	persistence.playerStateResolved(player)
	incorrectWord := game.Word("THISISNOTAWORD")
	nona.TryWord(player, incorrectWord)
	persistence.playerStateResolved(player)
	nona.GiveMe(player)
	persistence.playerStateResolved(player)
}

func TestPuzzles_CorrectSolution_UserIsNotified(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	puzzleMatch := puzzleSaver{}
	player := game.Player("U1")
	response := mock.NewMockResponse(mockCtrl)
	persistence := newFakePersistence()
	response.EXPECT().OnPuzzleNotification(player, &puzzleMatch)

	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player)
	persistence.playerStateResolved(player)
	correctWord := game.Word(oracle.FindASolutionFor(*puzzleMatch.puzzle))

	response.EXPECT().OnCorrectWord(player, correctWord)

	nona.TryWord(player, correctWord)
	persistence.playerStateResolved(player)
}

func TestPuzzles_IncorrectSolution_UserIsNotNotified(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	player := game.Player("U1")
	response := mock.NewMockResponse(mockCtrl)
	persistence := newFakePersistence()
	response.EXPECT().OnPuzzleNotification(player, gomock.Any())
	response.EXPECT().OnCorrectWord(player, gomock.Any()).Times(0)

	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player)
	persistence.playerStateResolved(player)
	incorrectWord := game.Word("THISISNOTAWORD")

	nona.TryWord(player, incorrectWord)
	persistence.playerStateResolved(player)
}

func TestTwoPlayers_FirstPuzzle_SameForBothPlayers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	firstPuzzleIsTheSame := samePuzzlesMatcher{}

	player1 := game.Player("U1")
	player2 := game.Player("U2")
	response := mock.NewMockResponse(mockCtrl)
	persistence := newFakePersistence()
	response.EXPECT().OnPuzzleNotification(player1, &firstPuzzleIsTheSame)
	response.EXPECT().OnPuzzleNotification(player2, &firstPuzzleIsTheSame)

	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player1)
	persistence.playerStateResolved(player1)
	nona.GiveMe(player2)
	persistence.playerStateResolved(player2)
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
	persistence := newFakePersistence()

	response.EXPECT().OnPuzzleNotification(player1, &player1Puzzle)
	response.EXPECT().OnPuzzleNotification(player2, &player2PuzzleUnchanged)
	response.EXPECT().OnPuzzleNotification(player2, &player2PuzzleUnchanged)
	response.EXPECT().OnCorrectWord(gomock.Any(), gomock.Any()).AnyTimes()

	// Arrange
	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player1)
	persistence.playerStateResolved(player1)
	nona.GiveMe(player2)
	persistence.playerStateResolved(player2)
	correctWord := game.Word(oracle.FindASolutionFor(*player1Puzzle.puzzle))
	nona.TryWord(player1, correctWord)
	persistence.playerStateResolved(player1)

	// Act
	// Here player1 has solved the first puzzle, but player2's puzzle should still be the first one.
	nona.GiveMe(player2)
	persistence.playerStateResolved(player2)
}
