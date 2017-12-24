package game_test

import (
	"testing"

	"github.com/dandeliondeathray/nona/game"
	"github.com/dandeliondeathray/nona/mock"
	"github.com/golang/mock/gomock"
)

func TestGamePersistence_NewRound_SeedIsStoredInPersistence(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	seed := int64(42)

	response := mock.NewMockResponse(mockCtrl)
	persistence := mock.NewMockPersistence(mockCtrl)
	persistence.EXPECT().StoreNewRound(seed)

	nona := game.NewGame(response, persistence, acceptanceDictionary)
	nona.NewRound(seed)
}
