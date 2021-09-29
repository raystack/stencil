package proto

import (
	"github.com/odpf/stencil/models"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Snapshot represents specific version of protodescriptorset
type Snapshot struct {
	ID        int64
	Namespace string
	Name      string
	Version   string
	Latest    bool
}

// toProtobufDBFiles creates DB compatible types
func toProtobufDBFiles(files *protoregistry.Files) []*models.ProtobufDBFile {
	var dbFiles []*models.ProtobufDBFile
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		f := ToProtobufDBFile(fd)
		dbFiles = append(dbFiles, f)
		return true
	})
	return dbFiles
}

// ToProtobufDBFile converts protoreflect.FileDescriptor type ProtobufDBFile
func ToProtobufDBFile(file protoreflect.FileDescriptor) *models.ProtobufDBFile {
	filefd := protodesc.ToFileDescriptorProto(file)
	data, _ := proto.MarshalOptions{Deterministic: true}.Marshal(filefd)
	return &models.ProtobufDBFile{
		Data: data,
		SearchData: &models.SearchData{
			Path:         file.Path(),
			Dependencies: getAllDependencies(file),
			Messages:     getMessageList(file),
			Package:      string(file.Package()),
			Fields:       getFieldList(file),
		},
	}
}

// fromByteArrayToFileDescriptorSet converts list of FileDescriptorProto []byte to FileDescriptorSet
func fromByteArrayToFileDescriptorSet(byteFiles [][]byte) ([]byte, error) {
	fds := &descriptorpb.FileDescriptorSet{}
	for _, byteFile := range byteFiles {
		fd := &descriptorpb.FileDescriptorProto{}
		proto.Unmarshal(byteFile, fd)
		fds.File = append(fds.File, fd)
	}
	data, err := proto.Marshal(fds)
	return data, err
}
