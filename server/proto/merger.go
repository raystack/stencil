package proto

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	defaultSkippedRules = []string{
		"FIELD_NO_BREAKING_CHANGE",
	}
)

func getFileDescriptorProto(fds *protoregistry.Files, path string) *descriptorpb.FileDescriptorProto {
	var fdp *descriptorpb.FileDescriptorProto
	fds.RangeFiles(func(fileDescriptor protoreflect.FileDescriptor) bool {
		if fileDescriptor.Path() == path {
			fdp = protodesc.ToFileDescriptorProto(fileDescriptor)
		}
		return true
	})
	return fdp
}

func getDescriptorProto(fdp *descriptorpb.FileDescriptorProto, name *string) *descriptorpb.DescriptorProto {
	for _, dp := range fdp.GetMessageType() {
		if *dp.Name == *name {
			return dp
		}
	}
	return nil
}

func validate(current []byte, prev []byte, rulesToSkip []string) error {
	for _, rule := range defaultSkippedRules {
		rulesToSkip = append(rulesToSkip, rule)
	}
	if err := Compare(current, prev, rulesToSkip); err != nil {
		return err
	}
	return nil
}

func Merge(current, prev []byte, rulesToSkip []string) ([]byte, error) {
	if err := validate(current, prev, rulesToSkip); err != nil {
		return nil, err
	}

	var err error
	var currentRegistry, previousRegistry *protoregistry.Files
	if currentRegistry, err = getRegistry(current); err != nil {
		return nil, err
	}
	if previousRegistry, err = getRegistry(prev); err != nil {
		return nil, err
	}

	var fileDescriptorProtos []*descriptorpb.FileDescriptorProto

	previousRegistry.RangeFiles(func(previousFD protoreflect.FileDescriptor) bool {
		previousFileDP := protodesc.ToFileDescriptorProto(previousFD)
		currentFileDP := getFileDescriptorProto(currentRegistry, previousFD.Path())

		// merge existing descriptor
		for _, previousDP := range previousFileDP.GetMessageType() {
			currentDP := getDescriptorProto(currentFileDP, previousDP.Name)
			mergeDescriptorProto(previousDP, currentDP)
		}

		addNewDependencies(currentFileDP, previousFileDP)
		addNewMessage(currentFileDP, previousFileDP)

		// add file package and options
		previousFileDP.Package = currentFileDP.Package
		previousFileDP.Options = currentFileDP.Options

		fileDescriptorProtos = append(fileDescriptorProtos, previousFileDP)
		return true
	})

	fileDescriptorProtos = addNewFiles(currentRegistry, previousRegistry, fileDescriptorProtos)

	return proto.
		MarshalOptions{Deterministic: true}.
		Marshal(&descriptorpb.FileDescriptorSet{
			File: fileDescriptorProtos,
		})
}

func addNewFiles(
	currentRegistry *protoregistry.Files,
	previousRegistry *protoregistry.Files,
	fileDescriptorProtos []*descriptorpb.FileDescriptorProto,
) []*descriptorpb.FileDescriptorProto {
	currentRegistry.RangeFiles(func(currentFD protoreflect.FileDescriptor) bool {
		currentFileDP := protodesc.ToFileDescriptorProto(currentFD)
		if getFileDescriptorProto(previousRegistry, currentFD.Path()) == nil {
			fileDescriptorProtos = append(fileDescriptorProtos, currentFileDP)
		}
		return true
	})
	return fileDescriptorProtos
}

func addNewMessage(currentFileDP *descriptorpb.FileDescriptorProto, previousFileDP *descriptorpb.FileDescriptorProto) {
	for _, currentDP := range currentFileDP.GetMessageType() {
		isNew := true
		for _, previousDP := range previousFileDP.GetMessageType() {
			if *previousDP.Name == *currentDP.Name {
				isNew = false
				break
			}
		}
		if isNew {
			previousFileDP.MessageType = append(previousFileDP.MessageType, currentDP)
		}
	}
}

func addNewDependencies(currentFileDP *descriptorpb.FileDescriptorProto, previousFileDP *descriptorpb.FileDescriptorProto) {
	// merge all new files and imports
	for _, currentDep := range currentFileDP.GetDependency() {
		isNew := true
		for _, previousDep := range previousFileDP.GetDependency() {
			if currentDep == previousDep {
				isNew = false
				break
			}
		}
		if isNew {
			previousFileDP.Dependency = append(previousFileDP.Dependency, currentDep)
		}
	}
}

func mergeDescriptorProto(previousDP *descriptorpb.DescriptorProto, currentDP *descriptorpb.DescriptorProto) {
	mergeNestedMessage(previousDP, currentDP)
	addNewNestedMessage(previousDP, currentDP)
	deprecateRemovedFields(previousDP, currentDP)
	addNewFields(previousDP, currentDP)
}

// merge from innermost nested message
func mergeNestedMessage(previousDP *descriptorpb.DescriptorProto, currentDP *descriptorpb.DescriptorProto) {
	for _, previousNestedDP := range previousDP.GetNestedType() {
		for _, currentNestedDP := range currentDP.GetNestedType() {
			if *previousNestedDP.Name == *currentNestedDP.Name {
				mergeDescriptorProto(previousNestedDP, currentNestedDP)
			}
		}
	}
}

func addNewNestedMessage(
	previousDP *descriptorpb.DescriptorProto,
	currentDP *descriptorpb.DescriptorProto,
) {
	for _, currentNestedDP := range currentDP.GetNestedType() {
		isNew := true
		for _, previousNestedDP := range previousDP.GetNestedType() {
			if *previousNestedDP.Name == *currentNestedDP.Name {
				isNew = false
				break
			}
		}
		if isNew {
			previousDP.NestedType = append(previousDP.NestedType, currentNestedDP)
		}
	}
}

func deprecateRemovedFields(previousDP *descriptorpb.DescriptorProto, currentDP *descriptorpb.DescriptorProto) {
	for _, previousFieldDP := range previousDP.Field {
		removed := true
		for _, currentFieldDP := range currentDP.Field {
			if *previousFieldDP.Name == *currentFieldDP.Name {
				removed = false
				break
			}
		}
		if removed {
			deprecated := true
			previousFieldDP.Options = &descriptorpb.FieldOptions{Deprecated: &deprecated}
		}
	}
}

func addNewFields(previousDP *descriptorpb.DescriptorProto, currentDP *descriptorpb.DescriptorProto) {
	// get max number in previous message
	var maxPreviousNumber int32 = 0
	for _, previousFieldDP := range previousDP.Field {
		if *previousFieldDP.Number > maxPreviousNumber {
			maxPreviousNumber = *previousFieldDP.Number
		}
	}
	// add new fields
	for _, currentFieldDP := range currentDP.Field {
		isNew := true
		for _, previousFieldDP := range previousDP.Field {
			if *currentFieldDP.Name == *previousFieldDP.Name {
				isNew = false
				break
			}
		}
		if isNew {
			fieldNumber := maxPreviousNumber + 1
			currentFieldDP.Number = &fieldNumber
			previousDP.Field = append(previousDP.Field, currentFieldDP)
			maxPreviousNumber = maxPreviousNumber + 1
		}
	}
}
