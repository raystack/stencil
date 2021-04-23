package stencil

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

type resolverBuilder struct {
	types *protoregistry.Types
}

func (r *resolverBuilder) addExtensions(exts protoreflect.ExtensionDescriptors) *resolverBuilder {
	for i := 0; i < exts.Len(); i++ {
		ext := exts.Get(i)
		r.types.RegisterExtension(dynamicpb.NewExtensionType(ext))
	}
	return r
}

func (r *resolverBuilder) addFromMessages(msgs protoreflect.MessageDescriptors) *resolverBuilder {
	for i := 0; i < msgs.Len(); i++ {
		msg := msgs.Get(i)
		r.addFromMessages(msg.Messages()).
			addExtensions(msg.Extensions())
	}
	return r
}

func (r *resolverBuilder) registerFile(file protoreflect.FileDescriptor) *resolverBuilder {
	return r.
		addExtensions(file.Extensions()).
		addFromMessages(file.Messages())
}

func getResolver(files *protoregistry.Files) *protoregistry.Types {
	types := &protoregistry.Types{}
	builder := &resolverBuilder{
		types: types,
	}
	files.RangeFiles(func(file protoreflect.FileDescriptor) bool {
		builder.registerFile(file)
		return true
	})
	return builder.types
}
