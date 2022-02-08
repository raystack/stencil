package stencil

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Resolver protobuf type resolver
type Resolver struct {
	types           *protoregistry.Types
	javaToProtoName map[string]string
}

// Get returns protobuf messageType for given proto message fullname.
// If java package file option is added, then message classname would be javapackage + message name.
// If java package file option is not defined then className would be proto message fullName.
func (r *Resolver) Get(className string) (protoreflect.MessageType, bool) {
	fullName, ok := r.javaToProtoName[className]
	if !ok {
		return nil, false
	}
	msg, err := r.types.FindMessageByName(protoreflect.FullName(fullName))
	if err != nil {
		return nil, false
	}
	return msg, true
}

// GetTypeResolver returns type resolver
func (r *Resolver) GetTypeResolver() *protoregistry.Types {
	return r.types
}

func (r *Resolver) addExtensions(exts protoreflect.ExtensionDescriptors) *Resolver {
	for i := 0; i < exts.Len(); i++ {
		ext := exts.Get(i)
		r.types.RegisterExtension(dynamicpb.NewExtensionType(ext))
	}
	return r
}

func (r *Resolver) addEnums(enums protoreflect.EnumDescriptors) *Resolver {
	for i := 0; i < enums.Len(); i++ {
		enum := enums.Get(i)
		r.types.RegisterEnum(dynamicpb.NewEnumType(enum))
	}
	return r
}

func (r *Resolver) addFromMessages(msgs protoreflect.MessageDescriptors) *Resolver {
	for i := 0; i < msgs.Len(); i++ {
		msg := msgs.Get(i)
		r.types.RegisterMessage(dynamicpb.NewMessageType(msg))
		r.javaToProtoName[defaultKeyFn(msg)] = string(msg.FullName())
		r.addFromMessages(msg.Messages()).
			addExtensions(msg.Extensions()).
			addEnums(msg.Enums())
	}
	return r
}

func (r *Resolver) registerFile(file protoreflect.FileDescriptor) *Resolver {
	return r.
		addExtensions(file.Extensions()).
		addEnums(file.Enums()).
		addFromMessages(file.Messages())
}

// NewResolver parses protobuf fileDescriptorSet schema returns type Resolver
func NewResolver(data []byte) (*Resolver, error) {
	types := &protoregistry.Types{}
	resolver := &Resolver{
		types:           types,
		javaToProtoName: make(map[string]string),
	}
	files, err := getFilesRegistry(data)
	if err != nil {
		return resolver, fmt.Errorf("file is not fully contained descriptor file. %w", err)
	}
	files.RangeFiles(func(file protoreflect.FileDescriptor) bool {
		resolver.registerFile(file)
		return true
	})
	return resolver, nil
}
