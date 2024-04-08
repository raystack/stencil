package changedetector

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/proto"
)

func IsMessageFieldChanged(field1, field2 *descriptor.FieldDescriptorProto) bool {
	return field1.GetType() != field2.GetType() || field1.GetName() != field2.GetName() || IsMessageFieldDeprecated(field1, field2)
}

func IsMessageFieldDeprecated(field1, field2 *descriptor.FieldDescriptorProto) bool {
	return field1.GetOptions() != field1.GetOptions() || field1.GetOptions().GetDeprecated() != field2.GetOptions().GetDeprecated()
}

func IsEnumFieldChanged(field1, field2 *descriptor.EnumValueDescriptorProto) bool {
	return field1.GetName() != field2.GetName() || IsEnumFieldDeprecated(field1, field2)
}

func IsEnumFieldDeprecated(field1, field2 *descriptor.EnumValueDescriptorProto) bool {
	return field1.GetOptions() != field1.GetOptions() || field1.GetOptions().GetDeprecated() != field2.GetOptions().GetDeprecated()
}

func GetDescriptorSet(data []byte) (*descriptor.FileDescriptorSet, error) {
	var fileDescriptorSet = &descriptor.FileDescriptorSet{}
	if err := proto.Unmarshal(data, fileDescriptorSet); err != nil {
		return nil, err
	}
	return fileDescriptorSet, nil
}

func GetImpactedMessageFields(oldMessageDesc, newMessageDesc *descriptor.DescriptorProto) []string {
	var impactedFields []string
	if oldMessageDesc.GetOptions().GetDeprecated() != newMessageDesc.GetOptions().GetDeprecated() {
		return append(impactedFields, oldMessageDesc.GetName())
	}
	var newFields = make(map[string]*descriptor.FieldDescriptorProto)
	for _, newField := range newMessageDesc.GetField() {
		newFields[newField.GetName()] = newField
	}
	for _, oldField := range oldMessageDesc.GetField() {
		if newFields[oldField.GetName()] != nil {
			if IsMessageFieldChanged(oldField, newFields[oldField.GetName()]) {
				impactedFields = append(impactedFields, oldField.GetName())
			}
			delete(newFields, oldField.GetName())
		}
	}
	for name := range newFields {
		impactedFields = append(impactedFields, name)
	}
	return append(impactedFields, GetImpactedEnumFieldInsideMessage(oldMessageDesc, newMessageDesc)...)
}

func GetImpactedEnumFieldInsideMessage(oldMessageDesc, newMessageDesc *descriptor.DescriptorProto) []string {
	var impactedEnums []string
	var newEnums = make(map[string]*descriptor.EnumDescriptorProto)
	for _, newEnum := range newMessageDesc.GetEnumType() {
		newEnums[newEnum.GetName()] = newEnum
	}
	for _, oldEnum := range oldMessageDesc.GetEnumType() {
		if newEnums[oldEnum.GetName()] != nil && !proto.Equal(oldEnum, newEnums[oldEnum.GetName()]) {
			impactedEnums = append(impactedEnums, oldEnum.GetName())
			delete(newEnums, oldEnum.GetName())
		}
	}
	for name := range newEnums {
		impactedEnums = append(impactedEnums, name)
	}
	return impactedEnums
}

func GetImpactedEnumFields(oldEnumDesc, newEnumDesc *descriptor.EnumDescriptorProto) []string {
	var impactedFields []string
	if oldEnumDesc.GetOptions().GetDeprecated() != newEnumDesc.GetOptions().GetDeprecated() {
		return append(impactedFields, oldEnumDesc.GetName())
	}
	var newFields = make(map[string]*descriptor.EnumValueDescriptorProto)
	for _, newField := range newEnumDesc.GetValue() {
		newFields[newField.GetName()] = newField
	}
	for _, oldField := range oldEnumDesc.GetValue() {
		if newFields[oldField.GetName()] != nil {
			if IsEnumFieldChanged(oldField, newFields[oldField.GetName()]) {
				impactedFields = append(impactedFields, oldEnumDesc.GetName()+"."+oldField.GetName())
			}
			delete(newFields, oldField.GetName())
		}
	}
	for name := range newFields {
		impactedFields = append(impactedFields, newEnumDesc.GetName()+"."+name)
	}
	return impactedFields
}
