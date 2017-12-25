// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dandeliondeathray/nona/slack (interfaces: Game)

// Package mock is a generated GoMock package.
package mock

import (
	game "github.com/dandeliondeathray/nona/game"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockGame is a mock of Game interface
type MockGame struct {
	ctrl     *gomock.Controller
	recorder *MockGameMockRecorder
}

// MockGameMockRecorder is the mock recorder for MockGame
type MockGameMockRecorder struct {
	mock *MockGame
}

// NewMockGame creates a new mock instance
func NewMockGame(ctrl *gomock.Controller) *MockGame {
	mock := &MockGame{ctrl: ctrl}
	mock.recorder = &MockGameMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGame) EXPECT() *MockGameMockRecorder {
	return m.recorder
}

// GiveMe mocks base method
func (m *MockGame) GiveMe(arg0 game.Player) {
	m.ctrl.Call(m, "GiveMe", arg0)
}

// GiveMe indicates an expected call of GiveMe
func (mr *MockGameMockRecorder) GiveMe(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GiveMe", reflect.TypeOf((*MockGame)(nil).GiveMe), arg0)
}

// TryWord mocks base method
func (m *MockGame) TryWord(arg0 game.Player, arg1 game.Word) {
	m.ctrl.Call(m, "TryWord", arg0, arg1)
}

// TryWord indicates an expected call of TryWord
func (mr *MockGameMockRecorder) TryWord(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryWord", reflect.TypeOf((*MockGame)(nil).TryWord), arg0, arg1)
}