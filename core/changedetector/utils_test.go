package changedetector_test

import (
	"testing"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/goto/stencil/core/changedetector"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestIsMessageFieldChanged(t *testing.T) {
	for _, test := range []struct {
		name     string
		field1   *descriptor.FieldDescriptorProto
		field2   *descriptor.FieldDescriptorProto
		expected bool
	}{
		{
			"should return true if message field type is changed",
			&descriptor.FieldDescriptorProto{
				Name: proto.String("field1"),
				Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
			},
			&descriptor.FieldDescriptorProto{
				Name: proto.String("field2"),
				Type: descriptor.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
			},
			true,
		},
		{
			"should return true if message field name is changed",
			&descriptor.FieldDescriptorProto{
				Name: proto.String("field1"),
				Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
			},
			&descriptor.FieldDescriptorProto{
				Name: proto.String("field2"),
				Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
			},
			true,
		},
		{
			"should return true if message field is deprecated",
			&descriptor.FieldDescriptorProto{
				Name: proto.String("field1"),
				Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
			},
			&descriptor.FieldDescriptorProto{
				Name: proto.String("field1"),
				Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
				Options: &descriptor.FieldOptions{
					Deprecated: proto.Bool(true),
				},
			},
			true,
		},
		{
			"should return false if message field is not changed",
			&descriptor.FieldDescriptorProto{
				Name: proto.String("field1"),
				Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
			},
			&descriptor.FieldDescriptorProto{
				Name: proto.String("field1"),
				Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
			},
			false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, changedetector.IsMessageFieldChanged(test.field1, test.field2))
		})
	}
}

func TestIsEnumFieldChanged(t *testing.T) {
	for _, test := range []struct {
		name     string
		field1   *descriptor.EnumValueDescriptorProto
		field2   *descriptor.EnumValueDescriptorProto
		expected bool
	}{
		{
			"should return true if enum field name is changed",
			&descriptor.EnumValueDescriptorProto{
				Name: proto.String("enum1"),
			},
			&descriptor.EnumValueDescriptorProto{
				Name: proto.String("enum2"),
			},
			true,
		},
		{
			"should return true if enum field is deprecated",
			&descriptor.EnumValueDescriptorProto{
				Name: proto.String("enum1"),
			},
			&descriptor.EnumValueDescriptorProto{
				Name: proto.String("enum2"),
				Options: &descriptor.EnumValueOptions{
					Deprecated: proto.Bool(true),
				},
			},
			true,
		},
		{
			"should return false if enum field is not changed",
			&descriptor.EnumValueDescriptorProto{
				Name: proto.String("enum1"),
			},
			&descriptor.EnumValueDescriptorProto{
				Name: proto.String("enum1"),
			},
			false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, changedetector.IsEnumFieldChanged(test.field1, test.field2))
		})
	}
}

func TestGetImpactedMessageFields(t *testing.T) {
	for _, test := range []struct {
		name     string
		field1   *descriptor.DescriptorProto
		field2   *descriptor.DescriptorProto
		expected []string
	}{
		{
			"should return impacted fields if any field got changed",
			&descriptor.DescriptorProto{
				Name: proto.String("user"),
				Field: []*descriptor.FieldDescriptorProto{
					{
						Name: proto.String("field1"),
						Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
				},
			},
			&descriptor.DescriptorProto{
				Name: proto.String("field1"),
				Field: []*descriptor.FieldDescriptorProto{
					{
						Name: proto.String("field1"),
						Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
						Options: &descriptor.FieldOptions{
							Deprecated: proto.Bool(true),
						},
					},
				},
			},
			[]string{"field1"},
		},
		{
			"should return impacted fields if any new field got added",
			&descriptor.DescriptorProto{
				Name: proto.String("user"),
				Field: []*descriptor.FieldDescriptorProto{
					{
						Name: proto.String("field1"),
						Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
				},
			},
			&descriptor.DescriptorProto{
				Name: proto.String("field1"),
				Field: []*descriptor.FieldDescriptorProto{
					{
						Name: proto.String("field1"),
						Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
					{
						Name: proto.String("field2"),
						Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
				},
			},
			[]string{"field2"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, changedetector.GetImpactedMessageFields(test.field1, test.field2))
		})
	}
}

func TestGetImpactedEnumFields(t *testing.T) {
	for _, test := range []struct {
		name     string
		field1   *descriptor.EnumDescriptorProto
		field2   *descriptor.EnumDescriptorProto
		expected []string
	}{
		{
			"should return impacted enum fields if any enum field got changed",
			&descriptor.EnumDescriptorProto{
				Name: proto.String("user"),
				Value: []*descriptor.EnumValueDescriptorProto{
					{
						Name: proto.String("admin"),
					},
				},
			},
			&descriptor.EnumDescriptorProto{
				Name: proto.String("user"),
				Value: []*descriptor.EnumValueDescriptorProto{
					{
						Name: proto.String("admin"),
						Options: &descriptor.EnumValueOptions{
							Deprecated: proto.Bool(true),
						},
					},
				},
			},
			[]string{"user.admin"},
		},
		{
			"should return impacted enum fields if any new enum field got added",
			&descriptor.EnumDescriptorProto{
				Name: proto.String("user"),
				Value: []*descriptor.EnumValueDescriptorProto{
					{
						Name: proto.String("admin"),
					},
				},
			},
			&descriptor.EnumDescriptorProto{
				Name: proto.String("user"),
				Value: []*descriptor.EnumValueDescriptorProto{
					{
						Name: proto.String("admin"),
					},
					{
						Name: proto.String("developer"),
					},
				},
			},
			[]string{"user.developer"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, changedetector.GetImpactedEnumFields(test.field1, test.field2))
		})
	}
}

func TestGetImpactedEnumFieldsInsideMessage(t *testing.T) {
	for _, test := range []struct {
		name     string
		field1   *descriptor.DescriptorProto
		field2   *descriptor.DescriptorProto
		expected []string
	}{
		{
			"should return impacted enum fields if any enum field inside message got changed",
			&descriptor.DescriptorProto{
				Name: proto.String("User"),
				EnumType: []*descriptor.EnumDescriptorProto{
					{
						Name: proto.String("type"),
						Value: []*descriptor.EnumValueDescriptorProto{
							{
								Name: proto.String("admin"),
							},
						},
					},
				},
			},
			&descriptor.DescriptorProto{
				Name: proto.String("User"),
				EnumType: []*descriptor.EnumDescriptorProto{
					{
						Name: proto.String("type"),
						Value: []*descriptor.EnumValueDescriptorProto{
							{
								Name: proto.String("admin"),
								Options: &descriptor.EnumValueOptions{
									Deprecated: proto.Bool(true),
								},
							},
						},
					},
				},
			},
			[]string{"type"},
		},
		{
			"should return impacted enum fields if any enum field inside message got added",
			&descriptor.DescriptorProto{
				Name: proto.String("User"),
				EnumType: []*descriptor.EnumDescriptorProto{
					{
						Name: proto.String("type"),
						Value: []*descriptor.EnumValueDescriptorProto{
							{
								Name: proto.String("admin"),
							},
						},
					},
				},
			},
			&descriptor.DescriptorProto{
				Name: proto.String("User"),
				EnumType: []*descriptor.EnumDescriptorProto{
					{
						Name: proto.String("type"),
						Value: []*descriptor.EnumValueDescriptorProto{
							{
								Name: proto.String("admin"),
							},
							{
								Name: proto.String("developer"),
							},
						},
					},
				},
			},
			[]string{"type"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, changedetector.GetImpactedEnumFieldInsideMessage(test.field1, test.field2))
		})
	}
}
