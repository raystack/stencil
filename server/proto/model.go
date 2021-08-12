package proto

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// SearchData contains searchable field information
type SearchData struct {
	Path         string   `json:"path"`
	Messages     []string `json:"messages"`
	Dependencies []string `json:"dependencies"`
}

// ProtobufDBFile structure to store for each file info in DB
type ProtobufDBFile struct {
	ID         int64
	SearchData *SearchData
	Data       []byte
}

// Snapshot represents specific version of protodescriptorset
type Snapshot struct {
	ID        int64
	Namespace string
	Name      string
	Version   string
	Latest    bool
}

// ToProtobufDBFiles creates DB compatible types
func ToProtobufDBFiles(files *protoregistry.Files) []*ProtobufDBFile {
	var dbFiles []*ProtobufDBFile
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		f := ToProtobufDBFile(fd)
		dbFiles = append(dbFiles, f)
		return true
	})
	return dbFiles
}

// ToProtobufDBFile converts protoreflect.FileDescriptor type ProtobufDBFile
func ToProtobufDBFile(file protoreflect.FileDescriptor) *ProtobufDBFile {
	filefd := protodesc.ToFileDescriptorProto(file)
	data, _ := proto.MarshalOptions{Deterministic: true}.Marshal(filefd)
	return &ProtobufDBFile{
		Data: data,
		SearchData: &SearchData{
			Path:         file.Path(),
			Dependencies: getAllDependencies(file),
			Messages:     getMessageList(file),
		},
	}
}

// FromByteArrayToFileDescriptorSet converts list of FileDescriptorProto []byte to FileDescriptorSet
func FromByteArrayToFileDescriptorSet(byteFiles [][]byte) ([]byte, error) {
	fds := &descriptorpb.FileDescriptorSet{}
	for _, byteFile := range byteFiles {
		fd := &descriptorpb.FileDescriptorProto{}
		proto.Unmarshal(byteFile, fd)
		fds.File = append(fds.File, fd)
	}
	data, err := proto.Marshal(fds)
	return data, err
}
