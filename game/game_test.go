package game_test

import (
	"strings"
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
	response.EXPECT().OnIncorrectWord(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

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

func TestPuzzles_CorrectSolutionButLowercase_UserIsNotified(t *testing.T) {
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
	correctWord := game.Word(strings.ToLower(oracle.FindASolutionFor(*puzzleMatch.puzzle)))

	response.EXPECT().OnCorrectWord(player, gomock.Any())

	nona.TryWord(player, correctWord)
	persistence.playerStateResolved(player)
}

func TestPuzzles_IncorrectSolution_UserIsToldItIsIncorrect(t *testing.T) {
	// Assert
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	player := game.Player("U1")
	incorrectWord := game.Word("THISISNOTAWORD")

	response := mock.NewMockResponse(mockCtrl)
	persistence := newFakePersistence()
	response.EXPECT().OnPuzzleNotification(player, gomock.Any())
	response.EXPECT().OnCorrectWord(player, gomock.Any()).Times(0)
	response.EXPECT().OnIncorrectWord(player, incorrectWord, gomock.Any(), gomock.Any())

	// Arrange
	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player)
	persistence.playerStateResolved(player)

	// Act
	nona.TryWord(player, incorrectWord)
	persistence.playerStateResolved(player)
}

func TestPuzzles_TryPuzzleAsSolutionInLowercase_NoMismatchingLetters(t *testing.T) {
	// Assert
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	player := game.Player("U1")

	puzzleMatch := puzzleSaver{}
	response := mock.NewMockResponse(mockCtrl)
	persistence := newFakePersistence()
	response.EXPECT().OnPuzzleNotification(player, &puzzleMatch)
	response.EXPECT().OnCorrectWord(player, gomock.Any()).Times(0)

	// We want to ensure that the letter matching is done in a case insensitive way.
	// Therefore, if we send in the puzzle in lowercase as a word, then (assuming the puzzle isn't)
	// also a correct word, there should not be a mismatch.
	response.EXPECT().OnIncorrectWord(player, gomock.Any(), "", "")

	// Arrange
	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player)
	persistence.playerStateResolved(player)

	// Here the player has been notified of the puzzle.
	lowercaseSameAsPuzzle := game.Word(strings.ToLower(string(*puzzleMatch.puzzle)))

	// Act
	nona.TryWord(player, lowercaseSameAsPuzzle)
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

func TestPuzzles_WordAndPuzzleMismatch_UserIsToldWhatLettersMismatched(t *testing.T) {
	// Assert
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	player := game.Player("U1")
	word := game.Word("ABCZXY")
	savedPuzzle := puzzleSaver{}

	response := mock.NewMockResponse(mockCtrl)
	persistence := newFakePersistence()
	response.EXPECT().OnPuzzleNotification(player, &savedPuzzle)
	response.EXPECT().OnCorrectWord(player, gomock.Any()).Times(0)

	// Arrange
	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.NewRound(0)
	nona.GiveMe(player)
	persistence.playerStateResolved(player)
	// Here the player should have been notified of the puzzle.
	// Calculate the difference between puzzle and word.
	tooMany, tooFew := game.Diff(string(word), string(*savedPuzzle.puzzle))
	response.EXPECT().OnIncorrectWord(player, gomock.Any(), tooMany, tooFew)

	// Act
	nona.TryWord(player, word)
	persistence.playerStateResolved(player)
}

//
// Round recovery
//
func TestNoRoundSet_GiveMeCommand_ErrorResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	player := game.Player("U1")
	response := mock.NewMockResponse(mockCtrl)
	persistence := newFakePersistence()

	// Assert
	response.EXPECT().OnNoRound(player)

	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.GiveMe(player)
}

func TestNoRoundSet_CheckSolution_ErrorResponse(t *testing.T) {
	// Arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	player := game.Player("U1")
	response := mock.NewMockResponse(mockCtrl)
	persistence := newFakePersistence()

	// Assert
	response.EXPECT().OnNoRound(player)

	// Act
	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.TryWord(player, game.Word("SOMEWORD"))
}

func TestRoundIsRecovered_GiveMeCommand_PuzzleReturned(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	seed := int64(42)
	player := game.Player("U1")
	response := mock.NewMockResponse(mockCtrl)
	persistence := newFakePersistence()

	// Assert
	response.EXPECT().OnPuzzleNotification(gomock.Any(), gomock.Any())

	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.OnRoundRecovered(seed)
	nona.GiveMe(player)
	persistence.playerStateResolved(player)
}
