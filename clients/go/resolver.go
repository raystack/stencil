package stencilclient

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

type resolverBuilder struct {
	types *protoregistry.Types
}

func (r *resolverBuilder) AddExtensions(exts protoreflect.ExtensionDescriptors) *resolverBuilder {
	for i := 0; i < exts.Len(); i++ {
		ext := exts.Get(i)
		r.types.RegisterExtension(dynamicpb.NewExtensionType(ext))
	}
	return r
}

func (r *resolverBuilder) AddFromMessages(msgs protoreflect.MessageDescriptors) *resolverBuilder {
	for i := 0; i < msgs.Len(); i++ {
		msg := msgs.Get(i)
		r.AddFromMessages(msg.Messages()).
			AddExtensions(msg.Extensions())
	}
	return r
}

func (r *resolverBuilder) RegisterFile(file protoreflect.FileDescriptor) *resolverBuilder {
	return r.
		AddExtensions(file.Extensions()).
		AddFromMessages(file.Messages())
}

func getResolver(files *protoregistry.Files) *protoregistry.Types {
	types := &protoregistry.Types{}
	builder := &resolverBuilder{
		types: types,
	}
	files.RangeFiles(func(file protoreflect.FileDescriptor) bool {
		builder.RegisterFile(file)
		return true
	})
	return builder.types
}
