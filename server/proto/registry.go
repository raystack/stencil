package proto

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var ErrNotFound = errors.New("not found")
var ErrCast = errors.New("not able to convert to desired type")

func getRegistry(data []byte) (*protoregistry.Files, error) {
	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(data, fds); err != nil {
		return nil, fmt.Errorf("descriptor set file is not valid. %w", err)
	}
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return files, fmt.Errorf("file is not fully contained descriptor file. hint: generate file descriptorset with --include_imports option. %w", err)
	}
	return files, err
}

func isMsgDescriptor(val protoreflect.Descriptor) (protoreflect.MessageDescriptor, bool) {
	value, ok := val.(protoreflect.MessageDescriptor)
	return value, ok
}

func isEnumDescriptor(val protoreflect.Descriptor) (protoreflect.EnumDescriptor, bool) {
	value, ok := val.(protoreflect.EnumDescriptor)
	return value, ok
}
func isFieldDescriptor(val protoreflect.Descriptor) (protoreflect.FieldDescriptor, bool) {
	value, ok := val.(protoreflect.FieldDescriptor)
	return value, ok
}

func getMsgDescriptorFromFiles(files *protoregistry.Files, name protoreflect.FullName) (protoreflect.MessageDescriptor, error) {
	val, err := files.FindDescriptorByName(name)
	if err != nil {
		return nil, err
	}
	msg, ok := isMsgDescriptor(val)
	if !ok {
		return nil, ErrCast
	}
	return msg, nil
}

func getEnumDescriptorFromFiles(files *protoregistry.Files, name protoreflect.FullName) (protoreflect.EnumDescriptor, error) {
	val, err := files.FindDescriptorByName(name)
	if err != nil {
		return nil, err
	}
	enum, ok := isEnumDescriptor(val)
	if !ok {
		return nil, ErrCast
	}
	return enum, nil
}

func getFieldDescriptorFromFiles(files *protoregistry.Files, name protoreflect.FullName) (protoreflect.FieldDescriptor, error) {
	val, err := files.FindDescriptorByName(name)
	if err != nil {
		return nil, err
	}
	field, ok := isFieldDescriptor(val)
	if !ok {
		return nil, ErrCast
	}
	return field, nil
}

func forEachMessage(msgs protoreflect.MessageDescriptors, f func(msg protoreflect.MessageDescriptor) bool) {
	for i := 0; i < msgs.Len(); i++ {
		msg := msgs.Get(i)
		if !f(msg) {
			break
		}
		forEachMessage(msg.Messages(), f)
	}
}

func forEachEnum(enums protoreflect.EnumDescriptors, f func(enum protoreflect.EnumDescriptor) bool) {
	for i := 0; i < enums.Len(); i++ {
		enum := enums.Get(i)
		if !f(enum) {
			break
		}
	}
}

func forEachEnumValues(enumValues protoreflect.EnumValueDescriptors, f func(protoreflect.EnumValueDescriptor) bool) {
	for i := 0; i < enumValues.Len(); i++ {
		if !f(enumValues.Get(i)) {
			break
		}
	}
}

func forEachField(fields protoreflect.FieldDescriptors, f func(enum protoreflect.FieldDescriptor) bool) {
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if !f(field) {
			break
		}
	}
}

func forEachEnumsFromFile(file protoreflect.FileDescriptor, f func(protoreflect.EnumDescriptor) bool) {
	forEachEnum(file.Enums(), f)
	forEachMessage(file.Messages(), func(msg protoreflect.MessageDescriptor) bool {
		forEachEnum(msg.Enums(), f)
		return true
	})
}

func forEachFieldsFromFile(file protoreflect.FileDescriptor, f func(protoreflect.FieldDescriptor) bool) {
	forEachMessage(file.Messages(), func(msg protoreflect.MessageDescriptor) bool {
		forEachField(msg.Fields(), f)
		return true
	})
}

func compareEnums(current, prev protoreflect.EnumDescriptor) error {
	var err error

	forEachEnumValues(prev.Values(), func(value protoreflect.EnumValueDescriptor) bool {
		fullName := value.FullName()
		name := fullName.Name()
		currentEnumValue := current.Values().ByName(name)
		if currentEnumValue == nil {
			err = multierr.Combine(err, fmt.Errorf("enumValue %s deleted from current version", value.FullName()))
			return true
		}
		if currentEnumValue.Number() != value.Number() {
			err = multierr.Combine(err, fmt.Errorf("enumValue %s number changed from %d to %d", value.FullName(), value.Number(), currentEnumValue.Number()))
			return true
		}
		return true
	})
	return err
}
