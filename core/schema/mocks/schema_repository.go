// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	schema "github.com/raystack/stencil/core/schema"
	mock "github.com/stretchr/testify/mock"
)

// SchemaRepository is an autogenerated mock type for the Repository type
type SchemaRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, namespace, _a2, metadata, versionID, schemaFile
func (_m *SchemaRepository) Create(ctx context.Context, namespace string, _a2 string, metadata *schema.Metadata, versionID string, schemaFile *schema.SchemaFile) (int32, error) {
	ret := _m.Called(ctx, namespace, _a2, metadata, versionID, schemaFile)

	var r0 int32
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *schema.Metadata, string, *schema.SchemaFile) int32); ok {
		r0 = rf(ctx, namespace, _a2, metadata, versionID, schemaFile)
	} else {
		r0 = ret.Get(0).(int32)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, *schema.Metadata, string, *schema.SchemaFile) error); ok {
		r1 = rf(ctx, namespace, _a2, metadata, versionID, schemaFile)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0, _a1, _a2
func (_m *SchemaRepository) Delete(_a0 context.Context, _a1 string, _a2 string) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteVersion provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *SchemaRepository) DeleteVersion(_a0 context.Context, _a1 string, _a2 string, _a3 int32) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int32) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *SchemaRepository) Get(_a0 context.Context, _a1 string, _a2 string, _a3 int32) ([]byte, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int32) []byte); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int32) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestVersion provides a mock function with given fields: _a0, _a1, _a2
func (_m *SchemaRepository) GetLatestVersion(_a0 context.Context, _a1 string, _a2 string) (int32, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 int32
	if rf, ok := ret.Get(0).(func(context.Context, string, string) int32); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(int32)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMetadata provides a mock function with given fields: _a0, _a1, _a2
func (_m *SchemaRepository) GetMetadata(_a0 context.Context, _a1 string, _a2 string) (*schema.Metadata, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *schema.Metadata
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *schema.Metadata); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*schema.Metadata)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: _a0, _a1
func (_m *SchemaRepository) List(_a0 context.Context, _a1 string) ([]schema.Schema, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []schema.Schema
	if rf, ok := ret.Get(0).(func(context.Context, string) []schema.Schema); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]schema.Schema)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListVersions provides a mock function with given fields: _a0, _a1, _a2
func (_m *SchemaRepository) ListVersions(_a0 context.Context, _a1 string, _a2 string) ([]int32, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []int32
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []int32); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int32)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateMetadata provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *SchemaRepository) UpdateMetadata(_a0 context.Context, _a1 string, _a2 string, _a3 *schema.Metadata) (*schema.Metadata, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *schema.Metadata
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *schema.Metadata) *schema.Metadata); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*schema.Metadata)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, *schema.Metadata) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSchemaRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewSchemaRepository creates a new instance of SchemaRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSchemaRepository(t mockConstructorTestingTNewSchemaRepository) *SchemaRepository {
	mock := &SchemaRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
