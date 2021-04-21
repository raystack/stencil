package stencil

import (
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type descriptorStore struct {
	sync.RWMutex
	data              map[string]protoreflect.MessageDescriptor
	extensionResolver *protoregistry.Types
}

func newStore() *descriptorStore {
	return &descriptorStore{data: make(map[string]protoreflect.MessageDescriptor)}
}

func (s *descriptorStore) get(key string) (protoreflect.MessageDescriptor, bool) {
	s.RLock()
	defer s.RUnlock()
	d, ok := s.data[key]
	return d, ok
}

func (s *descriptorStore) load(data []byte, getKeyFn func(protoreflect.FileDescriptor, protoreflect.MessageDescriptor) string) error {
	msg := &descriptorpb.FileDescriptorSet{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return fmt.Errorf("invalid file descriptorset file. %w", err)
	}
	files, err := protodesc.NewFiles(msg)
	if err != nil {
		return fmt.Errorf("file is not fully contained descriptor file.%w", err)
	}
	newData := make(map[string]protoreflect.MessageDescriptor)
	files.RangeFiles(func(file protoreflect.FileDescriptor) bool {
		forEachMessage(file.Messages(), func(msg protoreflect.MessageDescriptor) bool {
			key := getKeyFn(file, msg)
			newData[key] = msg
			return true
		})
		return true
	})
	resolver := getResolver(files)
	s.Lock()
	defer s.Unlock()
	s.extensionResolver = resolver
	s.data = newData
	return nil
}

func (s *descriptorStore) loadFromURI(uri string, options Options) error {
	data, err := downloader(uri, options.HTTPOptions)
	if err != nil {
		return err
	}
	err = s.load(data, defaultKeyFn)
	return err
}
