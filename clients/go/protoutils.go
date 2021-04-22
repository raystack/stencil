package stencil

import (
	"strings"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func forEachMessage(msgs protoreflect.MessageDescriptors, f func(protoreflect.MessageDescriptor)) {
	for i := 0; i < msgs.Len(); i++ {
		msg := msgs.Get(i)
		f(msg)
		forEachMessage(msg.Messages(), f)
	}
}

func getJavaPackage(fileDesc protoreflect.FileDescriptor) string {
	file := protodesc.ToFileDescriptorProto(fileDesc)
	options := file.Options
	if options != nil && options.JavaPackage != nil {
		return *options.JavaPackage
	}
	return ""
}

func defaultKeyFn(file protoreflect.FileDescriptor, msg protoreflect.MessageDescriptor) string {
	fullName := string(msg.FullName())
	protoPackage := string(file.Package())
	pkg := getJavaPackage(file)
	if pkg == "" {
		return fullName
	}
	return strings.Replace(fullName, protoPackage, pkg, 1)
}
