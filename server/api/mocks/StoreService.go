// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	snapshot "github.com/odpf/stencil/server/snapshot"
	mock "github.com/stretchr/testify/mock"
)

// StoreService is an autogenerated mock type for the StoreService type
type StoreService struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0, _a1, _a2
func (_m *StoreService) Get(_a0 context.Context, _a1 *snapshot.Snapshot, _a2 []string) ([]byte, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, *snapshot.Snapshot, []string) []byte); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *snapshot.Snapshot, []string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: _a0, _a1, _a2
func (_m *StoreService) Insert(_a0 context.Context, _a1 *snapshot.Snapshot, _a2 []byte) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *snapshot.Snapshot, []byte) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Validate provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *StoreService) Validate(_a0 context.Context, _a1 *snapshot.Snapshot, _a2 []byte, _a3 []string) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *snapshot.Snapshot, []byte, []string) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
