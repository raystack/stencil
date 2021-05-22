package proto

import (
	"fmt"

	"go.uber.org/multierr"
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
		forEachMessage(prevFile.Messages(), func(msg protoreflect.MessageDescriptor) bool {
			msgName := msg.FullName()
			_, notFoundErr := getMsgDescriptorFromFiles(current, msgName)
			if notFoundErr != nil {
				err = multierr.Combine(err, fmt.Errorf("%s has been removed in current version", msgName))
			}
			return true
		})
		return true
	})
	return err
}

func checkFieldNoBreakingChange(current, prev *protoregistry.Files) error {
	var err error
	prev.RangeFiles(func(prevFile protoreflect.FileDescriptor) bool {
		forEachFieldsFromFile(prevFile, func(prevField protoreflect.FieldDescriptor) bool {
			name := prevField.FullName()
			currentField, notFoundErr := getFieldDescriptorFromFiles(current, name)
			if notFoundErr != nil {
				err = multierr.Combine(err, fmt.Errorf("field %s is removed in current version", name))
				return true
			}
			if prevField == currentField {
				return true
			}
			if prevField.Kind() != currentField.Kind() {
				err = multierr.Combine(err, fmt.Errorf("type has changed for %s from %s to %s", name, prevField.Kind().String(), currentField.Kind().String()))
			}
			if prevField.Number() != currentField.Number() {
				err = multierr.Combine(err, fmt.Errorf("number changed for %s from %d to %d", name, prevField.Number(), currentField.Number()))
			}
			if prevField.Cardinality() != currentField.Cardinality() {
				err = multierr.Combine(err, fmt.Errorf("label changed for %s from %s to %s", name, prevField.Cardinality().String(), currentField.Cardinality().String()))
			}
			if prevField.JSONName() != currentField.JSONName() {
				err = multierr.Combine(err, fmt.Errorf("json name changed for %s from %s to %s", name, prevField.JSONName(), currentField.JSONName()))
			}
			return true
		})
		return true
	})
	return err
}

func checkFileNoBreakingChange(current, prev *protoregistry.Files) error {
	var err error
	prev.RangeFiles(func(prevFile protoreflect.FileDescriptor) bool {
		path := prevFile.Path()
		currentFile, notFoundErr := current.FindFileByPath(path)
		if notFoundErr != nil {
			err = multierr.Combine(err, fmt.Errorf("\"%s\" file has been deleted in current version", path))
			return true
		}
		if currentFile.Syntax() != prevFile.Syntax() {
			err = multierr.Combine(err, fmt.Errorf("syntax for %s changed from %s to %s", path, prevFile.Syntax().String(), currentFile.Syntax().String()))
		}
		if currentFile.Package() != prevFile.Package() {
			err = multierr.Combine(err, fmt.Errorf("package for %s changed from %s to %s", path, prevFile.Package(), currentFile.Package()))
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
			err = multierr.Combine(err, fmt.Errorf("all file options have been removed in %s current version", path))
			return true
		}
		if prevFileOptions.GetJavaPackage() != currentFileOptions.GetJavaPackage() {
			err = multierr.Combine(err, fmt.Errorf("java package for %s changed from %s to %s", path, prevFileOptions.GetJavaPackage(), currentFileOptions.GetJavaPackage()))
		}
		if prevFileOptions.GetJavaOuterClassname() != currentFileOptions.GetJavaOuterClassname() {
			err = multierr.Combine(err, fmt.Errorf("java outer classname for %s changed from %s to %s", path, prevFileOptions.GetJavaOuterClassname(), currentFileOptions.GetJavaOuterClassname()))
		}
		if prevFileOptions.GetGoPackage() != currentFileOptions.GetGoPackage() {
			err = multierr.Combine(err, fmt.Errorf("go package for %s changed from %s to %s", path, prevFileOptions.GetGoPackage(), currentFileOptions.GetGoPackage()))
		}
		return true
	})
	return err
}

func checkEnumNoBreakingChange(current, prev *protoregistry.Files) error {
	var err error
	prev.RangeFiles(func(prevFile protoreflect.FileDescriptor) bool {
		forEachEnumsFromFile(prevFile, func(prevEnum protoreflect.EnumDescriptor) bool {
			enumName := prevEnum.FullName()
			currentEnum, notFoundErr := getEnumDescriptorFromFiles(current, enumName)
			if notFoundErr != nil {
				err = multierr.Combine(err, fmt.Errorf("%s enum has been removed from current version", enumName))
				return true
			}
			checkErr := compareEnums(currentEnum, prevEnum)
			err = multierr.Combine(err, checkErr)
			return true
		})
		return true
	})
	return err
}
