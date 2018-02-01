package slack_test

import (
	"testing"

	"github.com/dandeliondeathray/nona/game"
	"github.com/dandeliondeathray/nona/mock"
	"github.com/dandeliondeathray/nona/slack"
	"github.com/golang/mock/gomock"
)

//
// Core commands
//

func TestCommand_GiveMe_CallGiveMeOnGame(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	g := mock.NewMockGame(mockCtrl)
	self := game.Player("USELF")
	player := game.Player("U1")

	g.EXPECT().GiveMe(player)

	handler := slack.NewNonaSlackHandler(g, self)
	handler.OnMessage(player, "!gemig")
}

func TestCommand_ASingleWord_CallTryWord(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	g := mock.NewMockGame(mockCtrl)
	self := game.Player("USELF")
	player := game.Player("U1")
	word := "PUSSGURKA"

	g.EXPECT().TryWord(player, game.Word(word))

	handler := slack.NewNonaSlackHandler(g, self)
	handler.OnMessage(player, word)
}

func TestCommand_Skippa_SkipPuzzle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	g := mock.NewMockGame(mockCtrl)
	self := game.Player("USELF")
	player := game.Player("U1")

	g.EXPECT().SkipPuzzle(player)

	handler := slack.NewNonaSlackHandler(g, self)
	handler.OnMessage(player, "!skippa")
}

//
// Trimming whitespace
//

func TestTrimming_TrailingSpaceInWord_SpaceIsTrimmed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	g := mock.NewMockGame(mockCtrl)
	self := game.Player("USELF")
	player := game.Player("U1")
	word := "PUSSGURKA"

	g.EXPECT().TryWord(player, game.Word(word))

	handler := slack.NewNonaSlackHandler(g, self)
	handler.OnMessage(player, word+" ")
}

func TestTrimming_PrefixSpaceInWord_SpaceIsTrimmed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	g := mock.NewMockGame(mockCtrl)
	self := game.Player("USELF")
	player := game.Player("U1")
	word := "PUSSGURKA"

	g.EXPECT().TryWord(player, game.Word(word))

	handler := slack.NewNonaSlackHandler(g, self)
	handler.OnMessage(player, " "+word)
}

func TestTrimming_PrefixSpaceInGiveMe_SpaceIsTrimmed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	g := mock.NewMockGame(mockCtrl)
	self := game.Player("USELF")
	player := game.Player("U1")

	g.EXPECT().GiveMe(player)

	handler := slack.NewNonaSlackHandler(g, self)
	handler.OnMessage(player, "  !gemig")
}

func TestTrimming_TrailingSpaceInGiveMe_SpaceIsTrimmed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	g := mock.NewMockGame(mockCtrl)
	self := game.Player("USELF")
	player := game.Player("U1")

	g.EXPECT().GiveMe(player)

	handler := slack.NewNonaSlackHandler(g, self)
	handler.OnMessage(player, "!gemig  ")
}

//
// Some commands have aliases for simplicity in writing on mobile devices.
//

func TestAlias_ExclamationIsAliasForGiveMe_CallGiveMe(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	g := mock.NewMockGame(mockCtrl)
	self := game.Player("USELF")
	player := game.Player("U1")

	g.EXPECT().GiveMe(player)

	handler := slack.NewNonaSlackHandler(g, self)
	handler.OnMessage(player, "!")
}

func TestAlias_SpaceInGiveMeForMobileSimplicity_CallGiveMe(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	g := mock.NewMockGame(mockCtrl)
	self := game.Player("USELF")
	player := game.Player("U1")

	g.EXPECT().GiveMe(player)

	handler := slack.NewNonaSlackHandler(g, self)
	handler.OnMessage(player, "! gemig")
}

//
// Slack special cases
//

func TestSelf_WordsFromSelf_Ignore(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	g := mock.NewMockGame(mockCtrl)
	self := game.Player("USELF")
	word := "PUSSGURKA"

	g.EXPECT().GiveMe(gomock.Any()).Times(0)
	g.EXPECT().TryWord(gomock.Any(), gomock.Any()).Times(0)

	handler := slack.NewNonaSlackHandler(g, self)
	handler.OnMessage(self, word)
}
