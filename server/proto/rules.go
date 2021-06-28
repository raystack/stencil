package proto

import (
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var (
	Rules = []Rule{
		NewRule(
			"FILE_NO_BREAKING_CHANGE",
			"Ensures file does not have breaking changes from previous version, eg: No change in syntax, package name or removal of file",
			checkFileNoBreakingChange,
		),
		NewRule(
			"MESSAGE_NO_DELETE",
			"Ensures no message is deleted from previous version",
			checkMessageNoDelete,
		),
		NewRule(
			"FIELD_NO_BREAKING_CHANGE",
			"Ensures field does not have breaking changes from previous version, eg: No change in type, number, json name, label, field name or removal of field",
			checkFieldNoBreakingChange,
		),
		NewRule(
			"ENUM_NO_BREAKING_CHANGE",
			"Ensures enum does not have breaking changes from previous version, eg: No change in name and number or removal of enum",
			checkEnumNoBreakingChange,
		),
	}
)

func checkMessageNoDelete(current, prev *protoregistry.Files) error {
	var err error
	prev.RangeFiles(func(prevFile protoreflect.FileDescriptor) bool {
		validationErr := newValidationErr(prevFile)
		forEachMessage(prevFile.Messages(), func(msg protoreflect.MessageDescriptor) bool {
			msgName := msg.FullName()
			_, notFoundErr := getMsgDescriptorFromFiles(current, msgName)
			if notFoundErr != nil {
				validationErr.add(`"%s" message has been removed`, msgName)
				return true
			}
			return true
		})
		err = combineErr(err, validationErr)
		return true
	})
	return err
}

func checkFieldNoBreakingChange(current, prev *protoregistry.Files) error {
	var err error
	checkMessagePair(current, prev, func(currentMsg, prevMsg protoreflect.MessageDescriptor) bool {
		validationErr := newValidationErr(prevMsg)
		prevFields := prevMsg.Fields()
		forEachField(prevFields, func(prevField protoreflect.FieldDescriptor) bool {
			name := prevField.FullName()

			currentField := currentMsg.Fields().ByName(prevField.Name())
			if currentField == nil {
				validationErr.add(`field "%s" is removed`, name)
				return true
			}
			if prevField == currentField {
				return true
			}
			if prevField.Kind() != currentField.Kind() {
				validationErr.add(`type has changed for "%s" from "%s" to "%s"`, name, prevField.Kind().String(), currentField.Kind().String())
			}
			if prevField.Kind() == currentField.Kind() {
				if prevField.Message() != nil && currentField.Message() != nil {
					prevMsgType := prevField.Message()
					currentMsgType := currentField.Message()
					if prevMsgType.FullName() != currentMsgType.FullName() {
						validationErr.add(`type has changed for "%s" from "%s" to "%s"`, name, prevMsgType.FullName(), currentMsgType.FullName())
					}
				}
				if prevField.Enum() != nil && currentField.Enum() != nil {
					prevEnumType := prevField.Enum()
					currentEnumType := currentField.Enum()
					if prevEnumType.FullName() != currentEnumType.FullName() {
						validationErr.add(`type has changed for "%s" from "%s" to "%s"`, name, prevEnumType.FullName(), currentEnumType.FullName())
					}
				}

			}
			if prevField.Number() != currentField.Number() {
				validationErr.add(`number changed for "%s" from "%d" to "%d"`, name, prevField.Number(), currentField.Number())
			}
			if prevField.Cardinality() != currentField.Cardinality() {
				validationErr.add(`label changed for "%s" from "%s" to "%s"`, name, prevField.Cardinality().String(), currentField.Cardinality().String())
			}
			if prevField.JSONName() != currentField.JSONName() {
				validationErr.add(`json name changed for "%s" from "%s" to "%s"`, name, prevField.JSONName(), currentField.JSONName())
			}
			return true
		})
		err = combineErr(err, validationErr)
		return true
	})
	return err
}

func compareFileOptions(err *validationErr, optionName, prev, current string) {
	if prev != current {
		err.add(`File option "%s" changed from "%s" to "%s"`, optionName, prev, current)
	}
}

func checkFileNoBreakingChange(current, prev *protoregistry.Files) error {
	var err error
	prev.RangeFiles(func(prevFile protoreflect.FileDescriptor) bool {
		validationErr := newValidationErr(prevFile)
		path := prevFile.Path()
		currentFile, notFoundErr := current.FindFileByPath(path)
		defer func() {
			err = combineErr(err, validationErr)
		}()
		if notFoundErr != nil {
			validationErr.add(`file has been deleted`)
			return true
		}
		if currentFile.Syntax() != prevFile.Syntax() {
			validationErr.add(`syntax changed from "%s" to "%s"`, prevFile.Syntax().String(), currentFile.Syntax().String())
		}
		if currentFile.Package() != prevFile.Package() {
			validationErr.add(`package changed from "%s" to "%s"`, prevFile.Package(), currentFile.Package())
		}
		prevFileOptions := protodesc.ToFileDescriptorProto(prevFile).Options
		currentFileOptions := protodesc.ToFileDescriptorProto(currentFile).Options
		if prevFileOptions == nil && currentFileOptions == nil {
			return true
		}
		if prevFileOptions == nil && currentFileOptions != nil {
			return true
		}
		if prevFileOptions != nil && currentFileOptions == nil {
			validationErr.add("all file options have been removed in current version")
			return true
		}
		compareFileOptions(validationErr, "java package", prevFileOptions.GetJavaPackage(), currentFileOptions.GetJavaPackage())
		compareFileOptions(validationErr, "java outer classname", prevFileOptions.GetJavaOuterClassname(), currentFileOptions.GetJavaOuterClassname())
		return true
	})
	return err
}

func compareEnums(validationErr *validationErr, current, prev protoreflect.EnumDescriptor) {
	forEachEnumValues(prev.Values(), func(value protoreflect.EnumValueDescriptor) bool {
		fullName := value.FullName()
		name := fullName.Name()
		currentEnumValue := current.Values().ByName(name)
		if currentEnumValue == nil {
			validationErr.add(`enumValue "%s" deleted from enum "%s"`, name, value.Parent().FullName())
			return true
		}
		if currentEnumValue.Number() != value.Number() {
			validationErr.add(`enumValue "%s" number changed from "%d" to "%d"`, value.FullName(), value.Number(), currentEnumValue.Number())
			return true
		}
		return true
	})
}

func checkEnumNoBreakingChange(current, prev *protoregistry.Files) error {
	var err error
	prev.RangeFiles(func(prevFile protoreflect.FileDescriptor) bool {
		validationErr := newValidationErr(prevFile)
		forEachEnumsFromFile(prevFile, func(prevEnum protoreflect.EnumDescriptor) bool {
			enumName := prevEnum.FullName()
			currentEnum, notFoundErr := getEnumDescriptorFromFiles(current, enumName)
			if notFoundErr != nil {
				validationErr.add(`enum "%s" has been removed`, enumName)
				return true
			}
			compareEnums(validationErr, currentEnum, prevEnum)
			return true
		})
		err = combineErr(err, validationErr)
		return true
	})
	return err
}
