package proto

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Merge merges current FileDescriptorSet to existing FileDescriptorSet
func Merge(current, prev []byte) ([]byte, error) {
	var err error
	var currentRegistry, previousRegistry *protoregistry.Files
	if currentRegistry, err = getRegistry(current); err != nil {
		return nil, err
	}
	if previousRegistry, err = getRegistry(prev); err != nil {
		return nil, err
	}
	var mergedFiles []*descriptorpb.FileDescriptorProto
	previousRegistry.RangeFiles(func(prevFileDesc protoreflect.FileDescriptor) bool {
		prevFile := protodesc.ToFileDescriptorProto(prevFileDesc)
		currFile := getFileDescriptorProto(currentRegistry, prevFileDesc.Path())
		if currFile == nil {
			return true
		}
		mergeExistingMessage(prevFile, prevFileDesc.Messages(), currentRegistry)
		addNewDependencies(prevFile, currFile)
		addNewMessage(prevFile, currFile)
		updatePackage(prevFile, currFile)
		updateOptions(prevFile, currFile)
		mergedFiles = append(mergedFiles, prevFile)
		return true
	})
	mergedFiles = addNewFiles(currentRegistry, previousRegistry, mergedFiles)
	return proto.Marshal(&descriptorpb.FileDescriptorSet{File: mergedFiles})
}

// mergeExistingMessage merges updated message to existing message with same name.
// Existing message definition in prevFile is modified with changes from updated message with same name.
func mergeExistingMessage(prevFile *descriptorpb.FileDescriptorProto, prevMsgDescriptors protoreflect.MessageDescriptors, currentRegistry *protoregistry.Files) {
	var parent *descriptorpb.DescriptorProto
	var parents []*descriptorpb.DescriptorProto
	forEachMessage(prevMsgDescriptors, func(prevMsgDesc protoreflect.MessageDescriptor) bool {
		// check if message is removed.
		// if message is removed, deprecate it.
		prev := protodesc.ToDescriptorProto(prevMsgDesc)
		currMsg, err := getMsgDescriptorFromFiles(currentRegistry, prevMsgDesc.FullName())
		if err != nil {
			deprecateMessage(prev)
			var parentMsgs []*descriptorpb.DescriptorProto
			if _, ok := isFileDescriptor(prevMsgDesc.Parent()); ok {
				parentMsgs = prevFile.MessageType
			}
			if _, ok := isMsgDescriptor(prevMsgDesc.Parent()); ok {
				parentMsgs = parent.NestedType
			}
			replaceMessage(parentMsgs, prev)
			parent = prev
			return true
		}
		curr := protodesc.ToDescriptorProto(currMsg)
		// core merge logic
		addNewNestedMessage(prev, curr)
		deprecateRemovedFields(prev, curr)
		addNewFields(prev, curr)
		// as prev is merged with curr, replace existing prev in parent
		if _, ok := isFileDescriptor(prevMsgDesc.Parent()); ok {
			replaceMessage(prevFile.MessageType, prev)
			parent = prev
		}
		if parent != prev {
			replaceMessage(parent.NestedType, prev)
		}
		// if message is the last nested message of its parent,
		// set parent pointer to parent's parent.
		if prevMsgDesc.Messages().Len() > 0 {
			parents = append(parents, parent)
			parent = prev
		} else {
			if md, ok := isMsgDescriptor(prevMsgDesc.Parent()); ok {
				lastIndex := md.Messages().Len() - 1
				if md.Messages().Get(lastIndex) == prevMsgDesc {
					parent = parents[len(parents)-1]
					parents[len(parents)-1] = nil
				}
			}
		}
		return true
	})
}

// Get FileDescriptor from a registry by path and convert it to FileDescriptorProto
func getFileDescriptorProto(registry *protoregistry.Files, path string) *descriptorpb.FileDescriptorProto {
	fileDescriptor, err := registry.FindFileByPath(path)
	if err != nil {
		return nil
	}
	return protodesc.ToFileDescriptorProto(fileDescriptor)
}

// add new files from current registry to previous registry, store them in a list
func addNewFiles(
	currRegistry *protoregistry.Files,
	prevRegistry *protoregistry.Files,
	mergedFiles []*descriptorpb.FileDescriptorProto,
) []*descriptorpb.FileDescriptorProto {
	currRegistry.RangeFiles(func(currFileDesc protoreflect.FileDescriptor) bool {
		currFile := protodesc.ToFileDescriptorProto(currFileDesc)
		if getFileDescriptorProto(prevRegistry, currFileDesc.Path()) == nil {
			mergedFiles = append(mergedFiles, currFile)
		}
		return true
	})
	return mergedFiles
}

// add new message from currMsgs to prevMsgs
func addMessage(prevMsgs *[]*descriptorpb.DescriptorProto, currMsgs *[]*descriptorpb.DescriptorProto) {
	for _, curr := range *currMsgs {
		exist := false
		for _, prev := range *prevMsgs {
			if prev.GetName() == curr.GetName() {
				exist = true
				break
			}
		}
		if !exist {
			*prevMsgs = append(*prevMsgs, curr)
		}
	}
}

// Append all new message in current file to previous file.
func addNewMessage(prevFile *descriptorpb.FileDescriptorProto, currFile *descriptorpb.FileDescriptorProto) {
	addMessage(&prevFile.MessageType, &currFile.MessageType)
}

// Append all new nested message in current message to previous message.
func addNewNestedMessage(prevParent *descriptorpb.DescriptorProto, currParent *descriptorpb.DescriptorProto) {
	addMessage(&prevParent.NestedType, &currParent.NestedType)
}

// Add all new import statements in current file to previous file.
func addNewDependencies(prevFile *descriptorpb.FileDescriptorProto, currFile *descriptorpb.FileDescriptorProto) {
	for _, currDep := range currFile.GetDependency() {
		isNew := true
		for _, prevDep := range prevFile.GetDependency() {
			if currDep == prevDep {
				isNew = false
				break
			}
		}
		if isNew {
			prevFile.Dependency = append(prevFile.Dependency, currDep)
		}
	}
}

// Set deprecated options to true for each removed field in message.
func deprecateRemovedFields(prev *descriptorpb.DescriptorProto, curr *descriptorpb.DescriptorProto) {
	for _, prevField := range prev.Field {
		removed := true
		for _, currField := range curr.Field {
			if prevField.GetName() == currField.GetName() {
				removed = false
				break
			}
		}
		if removed {
			deprecated := true
			prevField.Options = &descriptorpb.FieldOptions{Deprecated: &deprecated}
		}
	}
}

// New field is added after the max number of existing fields.
func addNewFields(prev *descriptorpb.DescriptorProto, curr *descriptorpb.DescriptorProto) {
	// get max number in previous message
	var maxPrevNumber int32 = 0
	for _, prevField := range prev.Field {
		if *prevField.Number > maxPrevNumber {
			maxPrevNumber = *prevField.Number
		}
	}
	// add new fields after max number
	for _, currField := range curr.Field {
		isNew := true
		for _, prevField := range prev.Field {
			if currField.GetName() == prevField.GetName() {
				isNew = false
				break
			}
		}
		if isNew {
			fieldNumber := maxPrevNumber + 1
			currField.Number = &fieldNumber
			prev.Field = append(prev.Field, currField)
			maxPrevNumber = maxPrevNumber + 1
		}
	}
}

// Set deprecated options of message to true.
func deprecateMessage(msg *descriptorpb.DescriptorProto) {
	deprecated := true
	msg.Options = &descriptorpb.MessageOptions{Deprecated: &deprecated}
}

// Replace a message in list of message with a new message with same name.
func replaceMessage(msgs []*descriptorpb.DescriptorProto, newMsg *descriptorpb.DescriptorProto) {
	for i, m := range msgs {
		if m.GetName() == newMsg.GetName() {
			msgs[i] = newMsg
		}
	}
}

// update options in existing file with options in new file
func updateOptions(prevFile, currFile *descriptorpb.FileDescriptorProto) {
	prevFile.Options = currFile.Options
}

// update package name in existing file with package name in new file
func updatePackage(prevFile, currFile *descriptorpb.FileDescriptorProto) {
	prevFile.Package = currFile.Package
}
