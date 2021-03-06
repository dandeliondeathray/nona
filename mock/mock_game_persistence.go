// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dandeliondeathray/nona/game (interfaces: Persistence)

package mock

import (
	game "github.com/dandeliondeathray/nona/game"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPersistence is a mock of Persistence interface
type MockPersistence struct {
	ctrl     *gomock.Controller
	recorder *MockPersistenceMockRecorder
}

// MockPersistenceMockRecorder is the mock recorder for MockPersistence
type MockPersistenceMockRecorder struct {
	mock *MockPersistence
}

// NewMockPersistence creates a new mock instance
func NewMockPersistence(ctrl *gomock.Controller) *MockPersistence {
	mock := &MockPersistence{ctrl: ctrl}
	mock.recorder = &MockPersistenceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPersistence) EXPECT() *MockPersistenceMockRecorder {
	return m.recorder
}

// PlayerSkippedPuzzle mocks base method
func (m *MockPersistence) PlayerSkippedPuzzle(arg0 game.Player, arg1, arg2 int) {
	m.ctrl.Call(m, "PlayerSkippedPuzzle", arg0, arg1, arg2)
}

// PlayerSkippedPuzzle indicates an expected call of PlayerSkippedPuzzle
func (mr *MockPersistenceMockRecorder) PlayerSkippedPuzzle(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlayerSkippedPuzzle", reflect.TypeOf((*MockPersistence)(nil).PlayerSkippedPuzzle), arg0, arg1, arg2)
}

// PlayerSolvedPuzzle mocks base method
func (m *MockPersistence) PlayerSolvedPuzzle(arg0 game.Player, arg1 int) {
	m.ctrl.Call(m, "PlayerSolvedPuzzle", arg0, arg1)
}

// PlayerSolvedPuzzle indicates an expected call of PlayerSolvedPuzzle
func (mr *MockPersistenceMockRecorder) PlayerSolvedPuzzle(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlayerSolvedPuzzle", reflect.TypeOf((*MockPersistence)(nil).PlayerSolvedPuzzle), arg0, arg1)
}

// ResolveAllPlayerStates mocks base method
func (m *MockPersistence) ResolveAllPlayerStates(arg0 game.AllPlayerStatesResolution) {
	m.ctrl.Call(m, "ResolveAllPlayerStates", arg0)
}

// ResolveAllPlayerStates indicates an expected call of ResolveAllPlayerStates
func (mr *MockPersistenceMockRecorder) ResolveAllPlayerStates(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveAllPlayerStates", reflect.TypeOf((*MockPersistence)(nil).ResolveAllPlayerStates), arg0)
}

// ResolvePlayerState mocks base method
func (m *MockPersistence) ResolvePlayerState(arg0 game.Player, arg1 game.PlayerStateResolution) {
	m.ctrl.Call(m, "ResolvePlayerState", arg0, arg1)
}

// ResolvePlayerState indicates an expected call of ResolvePlayerState
func (mr *MockPersistenceMockRecorder) ResolvePlayerState(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolvePlayerState", reflect.TypeOf((*MockPersistence)(nil).ResolvePlayerState), arg0, arg1)
}

// StoreNewRound mocks base method
func (m *MockPersistence) StoreNewRound(arg0 int64) {
	m.ctrl.Call(m, "StoreNewRound", arg0)
}

// StoreNewRound indicates an expected call of StoreNewRound
func (mr *MockPersistenceMockRecorder) StoreNewRound(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreNewRound", reflect.TypeOf((*MockPersistence)(nil).StoreNewRound), arg0)
}
