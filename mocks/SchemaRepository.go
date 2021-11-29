// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/odpf/stencil/domain"
	mock "github.com/stretchr/testify/mock"
)

// SchemaRepository is an autogenerated mock type for the SchemaRepository type
type SchemaRepository struct {
	mock.Mock
}

// CreateSchema provides a mock function with given fields: ctx, namespace, schema, metadata, versionID, schemaFile
func (_m *SchemaRepository) CreateSchema(ctx context.Context, namespace string, schema string, metadata *domain.Metadata, versionID string, schemaFile *domain.SchemaFile) (int32, error) {
	ret := _m.Called(ctx, namespace, schema, metadata, versionID, schemaFile)

	var r0 int32
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *domain.Metadata, string, *domain.SchemaFile) int32); ok {
		r0 = rf(ctx, namespace, schema, metadata, versionID, schemaFile)
	} else {
		r0 = ret.Get(0).(int32)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, *domain.Metadata, string, *domain.SchemaFile) error); ok {
		r1 = rf(ctx, namespace, schema, metadata, versionID, schemaFile)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSchema provides a mock function with given fields: _a0, _a1, _a2
func (_m *SchemaRepository) DeleteSchema(_a0 context.Context, _a1 string, _a2 string) error {
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

// GetLatestSchema provides a mock function with given fields: _a0, _a1, _a2
func (_m *SchemaRepository) GetLatestSchema(_a0 context.Context, _a1 string, _a2 string) ([]byte, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []byte); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
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

// GetSchema provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *SchemaRepository) GetSchema(_a0 context.Context, _a1 string, _a2 string, _a3 int32) ([]byte, error) {
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

// GetSchemaMetadata provides a mock function with given fields: _a0, _a1, _a2
func (_m *SchemaRepository) GetSchemaMetadata(_a0 context.Context, _a1 string, _a2 string) (*domain.Metadata, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *domain.Metadata
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *domain.Metadata); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Metadata)
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

// ListSchemas provides a mock function with given fields: _a0, _a1
func (_m *SchemaRepository) ListSchemas(_a0 context.Context, _a1 string) ([]string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
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

// UpdateSchemaMetadata provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *SchemaRepository) UpdateSchemaMetadata(_a0 context.Context, _a1 string, _a2 string, _a3 *domain.Metadata) (*domain.Metadata, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *domain.Metadata
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *domain.Metadata) *domain.Metadata); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Metadata)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, *domain.Metadata) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
