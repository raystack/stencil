// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	schema "github.com/odpf/stencil/server/schema"
	mock "github.com/stretchr/testify/mock"
)

// CompatibilityFn is an autogenerated mock type for the CompatibilityFn type
type CompatibilityFn struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, _a1
func (_m *CompatibilityFn) Execute(_a0 schema.ParsedSchema, _a1 []schema.ParsedSchema) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(schema.ParsedSchema, []schema.ParsedSchema) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
