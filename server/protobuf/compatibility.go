package protobuf

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type diffKind int

const (
	_ diffKind = iota
	messageDelete
	nonInclusivereservedRange
	nonInclusiceReservedNames
	fieldDelete
	fieldDeleteWithoutReservedNumber
	fieldDeleteWithoutReservedName
	fieldNameChange
	fieldLabelchange
	fieldKindChange
	fieldTypeChange
	enumDelete
	enumValueDelete
	enumValueDeleteWithoutReservedNumber
	enumValueDeleteWithoutReservedName
	enumValueNumberChange
	syntaxChange
)

var (
	backwardCompatibility = []diffKind{
		messageDelete,
		nonInclusivereservedRange,
		nonInclusiceReservedNames,
		fieldDelete,
		fieldNameChange,
		fieldLabelchange,
		fieldKindChange,
		fieldTypeChange,
		enumDelete,
		enumValueDelete,
		enumValueNumberChange,
		syntaxChange}
	forwardCompatibility = []diffKind{
		messageDelete,
		nonInclusivereservedRange,
		nonInclusiceReservedNames,
		fieldNameChange,
		fieldLabelchange,
		fieldKindChange,
		fieldTypeChange,
		fieldDeleteWithoutReservedNumber,
		fieldDeleteWithoutReservedName,
		enumDelete,
		enumValueDeleteWithoutReservedNumber,
		enumValueDeleteWithoutReservedName,
		enumValueNumberChange,
		syntaxChange}
	// fullCompatibility = []diffKind{
	// 	messageDelete,
	// 	nonInclusivereservedRange,
	// 	nonInclusiceReservedNames,
	// 	fieldNameChange,
	// 	fieldLabelchange,
	// 	fieldKindChange,
	// 	fieldTypeChange,
	// 	enumDelete,
	// 	enumValueDelete,
	// 	enumValueNumberChange,
	// 	syntaxChange}
)

func (d diffKind) contains(others []diffKind) bool {
	for _, v := range others {
		if v == d {
			return true
		}
	}
	return false
}

func compareSchemas(current, prev *protoregistry.Files, notAllowedChanges []diffKind) error {
	diffs := &compatibilityErr{notAllowed: notAllowedChanges}
	prev.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		forEachMessage(fd.Messages(), func(prevMsg protoreflect.MessageDescriptor) bool {
			currentMsg := getMessage(current, prevMsg.FullName())
			if currentMsg == nil {
				diffs.add(messageDelete, prevMsg, `"%s" is removed`, prevMsg.FullName())
				return true
			}
			compareSyntax(currentMsg, prevMsg, diffs)
			compareMessages(currentMsg, prevMsg, diffs)
			return true
		})
		forEachEnum(fd, func(ed protoreflect.EnumDescriptor) bool {
			currentEnum := getEnum(current, ed.FullName())
			if currentEnum == nil {
				diffs.add(enumDelete, ed, `enum "%s" deleted`, ed.FullName())
				return true
			}
			compareEnums(currentEnum, ed, diffs)
			return true
		})
		return true
	})
	if diffs.isEmpty() {
		return nil
	}
	return diffs
}

func compareMessages(current, prev protoreflect.MessageDescriptor, diffs *compatibilityErr) {
	prevRanges := prev.ReservedRanges()
	for i := 0; i < prevRanges.Len(); i++ {
		prevRange := prevRanges.Get(i)
		start, end := prevRange[0], prevRange[1]
		if !(current.ReservedRanges().Has(start) && current.ReservedRanges().Has(end-1)) {
			diffs.add(nonInclusivereservedRange, current, "previous reserved range (%d, %d) is not inclusive of current range", start, end)
		}
	}
	// check reserved names
	prevNames := prev.ReservedNames()
	for i := 0; i < prevNames.Len(); i++ {
		prevName := prevNames.Get(i)
		if !(current.ReservedNames().Has(prevName)) {
			diffs.add(nonInclusiceReservedNames, current, `previous reserved name "%s" is removed`, prevName)
		}
	}
	// check field compatibility
	prevFields := prev.Fields()
	for i := 0; i < prevFields.Len(); i++ {
		prevField := prevFields.Get(i)
		currentField := current.Fields().ByNumber(prevField.Number())
		if currentField == nil {
			diffs.add(fieldDelete, prev, `field "%s" is deleted`, prevField.Name())
			if !current.ReservedRanges().Has(prevField.Number()) {
				diffs.add(fieldDeleteWithoutReservedNumber, prev, `field "%s" with number "%d" is not marked as reserved after the delete`, prevField.Name(), prevField.Number())
			}
			if !current.ReservedNames().Has(prevField.Name()) {
				diffs.add(fieldDeleteWithoutReservedName, prev, `field "%s" not marked as reserved after the delete`, prevField.Name())
			}
			continue
		}
		compareFields(currentField, prevField, diffs)
	}
}

func compareFields(currentField, prevField protoreflect.FieldDescriptor, diffs *compatibilityErr) {
	if prevField.JSONName() != currentField.JSONName() {
		diffs.add(fieldNameChange, prevField, `JSON field name changed from "%s" to "%s"`, prevField.JSONName(), currentField.JSONName())
	}
	name := prevField.Name()
	if prevField.Cardinality().IsValid() && prevField.Cardinality().String() != currentField.Cardinality().String() {
		diffs.add(fieldLabelchange, prevField, `field "%s" label changed from "%s" to "%s"`, name, prevField.Cardinality().String(), currentField.Cardinality().String())
	}
	if prevField.Kind() != currentField.Kind() {
		diffs.add(fieldKindChange, prevField, `field "%s" kind changed from "%s" to "%s"`, name, prevField.Kind().String(), currentField.Kind().String())
	} else {
		if prevField.Kind() == protoreflect.MessageKind && prevField.Message().FullName() != currentField.Message().FullName() {
			diffs.add(fieldTypeChange, prevField, `field "%s" type changed from "%s" to "%s"`, name, prevField.Message().FullName(), currentField.Message().FullName())
		}
		if prevField.Kind() == protoreflect.EnumKind && prevField.Enum().FullName() != currentField.Enum().FullName() {
			diffs.add(fieldTypeChange, prevField, `field "%s" type changed from "%s" to "%s"`, name, prevField.Enum().FullName(), currentField.Enum().FullName())
		}
	}
}

func compareEnums(current, prev protoreflect.EnumDescriptor, diffs *compatibilityErr) {
	// check reserved numbers
	prevRanges := prev.ReservedRanges()
	for i := 0; i < prevRanges.Len(); i++ {
		prevRange := prevRanges.Get(i)
		start, end := prevRange[0], prevRange[1]
		if !(current.ReservedRanges().Has(start) && current.ReservedRanges().Has(end)) {
			if start == end {
				diffs.add(nonInclusivereservedRange, current, "previous reserved number (%d) is not inclusive of current range", start)
			} else {
				diffs.add(nonInclusivereservedRange, current, "previous reserved range (%d, %d) is not inclusive of current range", start, end)
			}
		}
	}
	// check reserved names
	prevNames := prev.ReservedNames()
	for i := 0; i < prevNames.Len(); i++ {
		prevName := prevNames.Get(i)
		if !(current.ReservedNames().Has(prevName)) {
			diffs.add(nonInclusiceReservedNames, current, `previous reserved name "%s" is removed`, prevName)
		}
	}
	// check enum values
	prevValues := prev.Values()
	for i := 0; i < prevValues.Len(); i++ {
		prevValue := prevValues.Get(i)
		currentValue := current.Values().ByName(prevValue.Name())
		if currentValue == nil {
			diffs.add(enumValueDelete, prev, `enum value "%s" with number "%d" is deleted from "%s"`, prevValue.Name(), prevValue.Number(), prev.FullName())
			if !current.ReservedRanges().Has(prevValue.Number()) {
				diffs.add(enumValueDeleteWithoutReservedNumber, prev, `enum value "%s" with number "%d" is not marked as reserved after the delete`, prevValue.Name(), prevValue.Number())
			}
			if !current.ReservedNames().Has(prevValue.Name()) {
				diffs.add(enumValueDeleteWithoutReservedName, prev, `enum value "%s" name not marked as reserved after the delete`, prevValue.Name())
			}
			continue
		}
		if prevValue.Number() != currentValue.Number() {
			diffs.add(enumValueNumberChange, prev, `enum value number for "%s" changed from "%d" to "%d"`, prevValue.FullName(), prevValue.Number(), currentValue.Number())
		}
	}
}

func compareSyntax(current, prev protoreflect.Descriptor, diffs *compatibilityErr) {
	if current.ParentFile().Syntax() != prev.Parent().Syntax() {
		diffs.add(syntaxChange, current, `syntax changed from "%s" to "%s"`, prev.ParentFile().Syntax(), current.ParentFile().Syntax())
	}
}
