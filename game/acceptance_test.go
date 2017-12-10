package game_test

import (
	"testing"

	"github.com/dandeliondeathray/nona/game"
	"github.com/dandeliondeathray/nona/mock"
	"github.com/golang/mock/gomock"
)

var dictionary = []string{"PUSSGURKA"}

func TestGiveMeCommand_ForANewRound_PuzzleIsReturned(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	player := game.Player("U1")
	response := mock.NewMockResponse(mockCtrl)
	response.EXPECT().OnPuzzleNotification(player, gomock.Any())

	nona := game.NewGame(response, dictionary)
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

	nona := game.NewGame(response, dictionary)
	nona.NewRound(0)
	nona.GiveMe(player)
	nona.TryWord(player, game.Word("PUSSGURKA"))
	nona.GiveMe(player)
}
