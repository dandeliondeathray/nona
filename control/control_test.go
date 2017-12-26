package control_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dandeliondeathray/nona/control"
	"github.com/dandeliondeathray/nona/mock"
	"github.com/golang/mock/gomock"
)

func TestNewRound_SeedIs42_GameGetsSeed42(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Assert
	target := mock.NewMockTarget(mockCtrl)
	target.EXPECT().NewRound(int64(42))

	// Arrange
	request := httptest.NewRequest("POST", "/round/42", nil)
	recorder := httptest.NewRecorder()

	router := control.NewRouter(target)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("HTTP response code was %d but expected 200 OK", recorder.Code)
	}
}

func TestNewRound_SeedIs17_GameGetsSeed17(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Assert
	target := mock.NewMockTarget(mockCtrl)
	target.EXPECT().NewRound(int64(17))

	// Arrange
	request := httptest.NewRequest("POST", "/round/17", nil)
	recorder := httptest.NewRecorder()

	router := control.NewRouter(target)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("HTTP response code was %d but expected 200 OK", recorder.Code)
	}
}
