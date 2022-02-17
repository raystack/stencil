package protobuf

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func forEachMessage(msgs protoreflect.MessageDescriptors, f func(msg protoreflect.MessageDescriptor) bool) {
	for i := 0; i < msgs.Len(); i++ {
		msg := msgs.Get(i)
		if !f(msg) {
			break
		}
		forEachMessage(msg.Messages(), f)
	}
}

func eachEnum(enums protoreflect.EnumDescriptors, f func(protoreflect.EnumDescriptor) bool) {
	for i := 0; i < enums.Len(); i++ {
		e := enums.Get(i)
		if !f(e) {
			return
		}
	}
}

func forEachEnum(fd protoreflect.FileDescriptor, f func(protoreflect.EnumDescriptor) bool) {
	eachEnum(fd.Enums(), f)
	forEachMessage(fd.Messages(), func(msg protoreflect.MessageDescriptor) bool {
		eachEnum(msg.Enums(), f)
		return true
	})
}

func forEachField(fields protoreflect.FieldDescriptors, f func(protoreflect.FieldDescriptor) bool) {
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if !f(field) {
			break
		}
	}
}

func getMessageFromFile(file protoreflect.FileDescriptor) []string {
	var messages []string
	forEachMessage(file.Messages(), func(msg protoreflect.MessageDescriptor) bool {
		messages = append(messages, string(msg.FullName()))
		return true
	})
	return messages
}

func getFieldsFromFile(file protoreflect.FileDescriptor) []string {
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

func getAllMessages(s *protoregistry.Files) []string {
	var allMessages []string
	s.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		messages := getMessageFromFile(fd)
		allMessages = append(allMessages, messages...)
		return true
	})
	return allMessages
}

func getAllFields(s *protoregistry.Files) []string {
	var fieldNames []string
	s.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		names := getFieldsFromFile(fd)
		fieldNames = append(fieldNames, names...)
		return true
	})
	return fieldNames
}

func getMessage(files *protoregistry.Files, fullName protoreflect.FullName) protoreflect.MessageDescriptor {
	desc, err := files.FindDescriptorByName(fullName)
	if err != nil {
		return nil
	}
	msg, ok := desc.(protoreflect.MessageDescriptor)
	if ok {
		return msg
	}
	return nil
}

func getEnum(files *protoregistry.Files, fullName protoreflect.FullName) protoreflect.EnumDescriptor {
	desc, err := files.FindDescriptorByName(fullName)
	if err != nil {
		return nil
	}
	enum, ok := desc.(protoreflect.EnumDescriptor)
	if ok {
		return enum
	}
	return nil
}
