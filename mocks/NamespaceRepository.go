// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/odpf/stencil/server/domain"
	mock "github.com/stretchr/testify/mock"
)

// NamespaceRepository is an autogenerated mock type for the NamespaceRepository type
type NamespaceRepository struct {
	mock.Mock
}

// CreateNamespace provides a mock function with given fields: _a0, _a1
func (_m *NamespaceRepository) CreateNamespace(_a0 context.Context, _a1 domain.Namespace) (domain.Namespace, error) {
	ret := _m.Called(_a0, _a1)

	var r0 domain.Namespace
	if rf, ok := ret.Get(0).(func(context.Context, domain.Namespace) domain.Namespace); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(domain.Namespace)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.Namespace) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteNamespace provides a mock function with given fields: _a0, _a1
func (_m *NamespaceRepository) DeleteNamespace(_a0 context.Context, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetNamespace provides a mock function with given fields: _a0, _a1
func (_m *NamespaceRepository) GetNamespace(_a0 context.Context, _a1 string) (domain.Namespace, error) {
	ret := _m.Called(_a0, _a1)

	var r0 domain.Namespace
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Namespace); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(domain.Namespace)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListNamespaces provides a mock function with given fields: _a0
func (_m *NamespaceRepository) ListNamespaces(_a0 context.Context) ([]string, error) {
	ret := _m.Called(_a0)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateNamespace provides a mock function with given fields: _a0, _a1
func (_m *NamespaceRepository) UpdateNamespace(_a0 context.Context, _a1 domain.Namespace) (domain.Namespace, error) {
	ret := _m.Called(_a0, _a1)

	var r0 domain.Namespace
	if rf, ok := ret.Get(0).(func(context.Context, domain.Namespace) domain.Namespace); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(domain.Namespace)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.Namespace) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
