package stencil

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

func getJavaPackage(fileDesc protoreflect.FileDescriptor) string {
	file := protodesc.ToFileDescriptorProto(fileDesc)
	options := file.Options
	if options != nil && options.JavaPackage != nil {
		return *options.JavaPackage
	}
	return ""
}

func defaultKeyFn(msg protoreflect.MessageDescriptor) string {
	fullName := string(msg.FullName())
	file := msg.ParentFile()
	protoPackage := string(file.Package())
	pkg := getJavaPackage(file)
	if pkg == "" {
		return fullName
	}
	return strings.Replace(fullName, protoPackage, pkg, 1)
}

func getFilesRegistry(data []byte) (*protoregistry.Files, error) {
	msg := &descriptorpb.FileDescriptorSet{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, fmt.Errorf("invalid file descriptorset file. %w", err)
	}
	return protodesc.NewFiles(msg)
}
