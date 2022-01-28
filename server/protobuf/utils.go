package protobuf

import (
	"strings"

	"github.com/jdkato/prose/v2"
	"github.com/thoas/go-funk"
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

func getMessageFromFile(file protoreflect.FileDescriptor) ([]string, []string) {
	var messages []string
	var messageSearchKeys []string
	forEachMessage(file.Messages(), func(msg protoreflect.MessageDescriptor) bool {
		messages = append(messages, string(msg.FullName()))
		messageSearchKeys = append(messageSearchKeys, getSearchKeys(msg)...)
		return true
	})
	return messages, messageSearchKeys
}

func getFieldsFromFile(file protoreflect.FileDescriptor) ([]string, []string) {
	var fields []string
	var fieldSearchKeys []string
	forEachMessage(file.Messages(), func(msg protoreflect.MessageDescriptor) bool {
		forEachField(msg.Fields(), func(fd protoreflect.FieldDescriptor) bool {
			fields = append(fields, string(fd.FullName()))
			fieldSearchKeys = append(fieldSearchKeys, getSearchKeys(fd)...)
			return true
		})
		return true
	})
	return fields, fieldSearchKeys
}

func getAllMessages(s *protoregistry.Files) ([]string, []string) {
	var allMessages []string
	var allMessageSearchKeys []string
	s.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		messages, messageSearchKeys := getMessageFromFile(fd)
		allMessages = append(allMessages, messages...)
		allMessageSearchKeys = append(allMessageSearchKeys, messageSearchKeys...)
		return true
	})
	return allMessages, allMessageSearchKeys
}

func getAllFields(s *protoregistry.Files) ([]string, []string) {
	var fieldNames []string
	var fieldSearchKeys []string
	s.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		names, searchKeys := getFieldsFromFile(fd)
		fieldNames = append(fieldNames, names...)
		fieldSearchKeys = append(fieldSearchKeys, searchKeys...)
		return true
	})
	return fieldNames, fieldSearchKeys
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

func getSearchKeys(d protoreflect.Descriptor) []string {
	file := d.ParentFile()
	file.Services()
	sls := file.SourceLocations()
	if sls.Len() == 0 {
		return []string{}
	}
	sl := sls.ByDescriptor(d)
	acceptedTags := []string{"NNP", "NNPS"}
	doc, err := prose.NewDocument(sl.LeadingComments, prose.WithExtraction(false), prose.WithSegmentation(false))
	tokens := make([]string, 0)
	if err != nil {
		return []string{}
	}
	var prefix string
	m := file.Messages().ByName(d.Name())
	if m != nil {
		prefix = "m"
	} else {
		prefix = "f"
	}

	for _, t := range doc.Tokens() {
		text := strings.ToLower(t.Text)
		if funk.Contains(acceptedTags, t.Tag) {
			if !funk.Contains(tokens, text) && text != strings.ToLower(string(d.Name())) {
				tokens = append(tokens, prefix+"_"+string(d.FullName())+"."+strings.ToLower(t.Text))
			}
		}
	}
	return tokens
}
