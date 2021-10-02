package proto

import (
	"errors"
	"fmt"

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

func isFileDescriptor(val protoreflect.Descriptor) (protoreflect.FileDescriptor, bool) {
	value, ok := val.(protoreflect.FileDescriptor)
	return value, ok
}

func isMsgDescriptor(val protoreflect.Descriptor) (protoreflect.MessageDescriptor, bool) {
	value, ok := val.(protoreflect.MessageDescriptor)
	return value, ok
}

func isEnumDescriptor(val protoreflect.Descriptor) (protoreflect.EnumDescriptor, bool) {
	value, ok := val.(protoreflect.EnumDescriptor)
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

func forEachField(fields protoreflect.FieldDescriptors, f func(protoreflect.FieldDescriptor) bool) {
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

func checkMessagePair(current, prev *protoregistry.Files, f func(protoreflect.MessageDescriptor, protoreflect.MessageDescriptor) bool) {
	prev.RangeFiles(func(prevFile protoreflect.FileDescriptor) bool {
		forEachMessage(prevFile.Messages(), func(prevMsg protoreflect.MessageDescriptor) bool {
			msgName := prevMsg.FullName()
			currentMsg, err := getMsgDescriptorFromFiles(current, msgName)
			if err != nil {
				return true
			}
			return f(currentMsg, prevMsg)
		})
		return true
	})
}

func getFileDescriptor(desc protoreflect.Descriptor) protoreflect.FileDescriptor {
	f := desc.ParentFile()
	if f != nil && f.Parent() == nil {
		return f
	}
	if desc.Parent() == nil {
		f, ok := isFileDescriptor(desc)
		if ok {
			return f
		}
		return getFileDescriptor(desc.Parent())
	}
	return getFileDescriptor(desc.Parent())
}

func getMessageList(file protoreflect.FileDescriptor) []string {
	var messages []string
	forEachMessage(file.Messages(), func(msg protoreflect.MessageDescriptor) bool {
		messages = append(messages, string(msg.FullName()))
		return true
	})
	return messages
}

func getFieldList(file protoreflect.FileDescriptor) []string {
	var fields []string
	forEachMessage(file.Messages(), func(msg protoreflect.MessageDescriptor) bool {
		forEachField(msg.Fields(), func(fd protoreflect.FieldDescriptor) bool {
			fields = append(fields, string(fd.FullName()))
			return true
		})
		return true
	})
	return fields
}

func getAllDependencies(file protoreflect.FileDescriptor) []string {
	var fileImports []string
	fileImports = append(fileImports, file.Path())
	for i := 0; i < file.Imports().Len(); i++ {
		imp := file.Imports().Get(i)
		dependentImports := getAllDependencies(imp)
		fileImports = append(fileImports, dependentImports...)
	}
	return fileImports
}
