package protobuf

import (
	"fmt"

	"github.com/goto/stencil/core/schema"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

func getRegistry(fds *descriptorpb.FileDescriptorSet, data []byte) (*protoregistry.Files, error) {
	if err := proto.Unmarshal(data, fds); err != nil {
		return nil, fmt.Errorf("descriptor set file is not valid. %w", err)
	}
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return files, fmt.Errorf("file is not fully contained descriptor file. hint: generate file descriptorset with --include_imports option. %w", err)
	}
	return files, err
}

// GetParsedSchema converts data into enriched data type to deal with protobuf schema
func GetParsedSchema(data []byte) (schema.ParsedSchema, error) {
	fds := &descriptorpb.FileDescriptorSet{}
	files, err := getRegistry(fds, data)
	if err != nil {
		return &Schema{
			isValid: false,
		}, err
	}
	orderedData, _ := proto.MarshalOptions{Deterministic: true}.Marshal(fds)
	return &Schema{
		isValid: err == nil,
		Files:   files,
		data:    orderedData,
	}, nil
}
