// Package stencil helps to download and refresh protobuf descriptors from remote server
// and provides helper functions to get protobuf schema descriptors and can parse the messages
// dynamically.
package stencil

import (
	"encoding/json"
	"io"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

var (
	//ErrNotFound default sentinel error if proto not found
	ErrNotFound = errors.New("not found")
	//ErrNotFound is for when descriptor does not match the message
	ErrInvalidDescriptor = errors.New("invalid descriptor")
)

// Client provides utility functions to parse protobuf messages at runtime.
// protobuf messages can be identified by specifying fully qualified generated proto java class name.
type Client interface {
	// Parse parses protobuf message from wire format to protoreflect.ProtoMessage given fully qualified name of proto message.
	// Returns ErrNotFound error if given class name is not found
	Parse(string, []byte) (protoreflect.ProtoMessage, error)
	// ParseWithRefresh parses protobuf message from wire format to `protoreflect.ProtoMessage` given fully qualified name of proto message.
	// Refreshes proto definitions if parsed message has unknown fields and parses the message again.
	// Returns ErrNotFound error if given class name is not found.
	ParseWithRefresh(string, []byte) (protoreflect.ProtoMessage, error)
	// Serialize serializes data to bytes given fully qualified name of proto message.
	// Returns ErrNotFound error if given class name is not found
	Serialize(string, interface{}) ([]byte, error)
	// Serialize serializes data to bytes given fully qualified name of proto message.
	// Refreshes proto definitions if parsed message has unknown fields and parses the message again.
	// Returns ErrNotFound error if given class name is not found.
	SerializeWithRefresh(string, interface{}) ([]byte, error)
	// GetDescriptor returns protoreflect.MessageDescriptor given fully qualified proto java class name
	GetDescriptor(string) (protoreflect.MessageDescriptor, error)
	// Close stops background refresh if configured.
	Close()
	// Refresh downloads latest proto definitions
	Refresh() error
}

// HTTPOptions options for http client
type HTTPOptions struct {
	// Timeout specifies a time limit for requests made by this client. Default to 10s.
	// `0` duration not allowed. Client will set to default value (i.e. 10s).
	Timeout time.Duration
	// Headers provide extra headers to be added in requests made by this client
	Headers map[string]string
}

// Options options for stencil client
type Options struct {
	// AutoRefresh boolean to enable or disable autorefresh. Default to false
	AutoRefresh bool
	// RefreshInterval refresh interval to fetch descriptor file from server. Default to 12h.
	// `0` duration not allowed. Client will set to default value (i.e. 12h).
	RefreshInterval time.Duration
	// HTTPOptions options for http client
	HTTPOptions
}

func (o *Options) setDefaults() {
	if o.RefreshInterval == 0 {
		o.RefreshInterval = 12 * time.Hour
	}
	if o.HTTPOptions.Timeout == 0 {
		o.HTTPOptions.Timeout = 10 * time.Second
	}
}

type stencilClient struct {
	timer   io.Closer
	urls    []string
	store   *descriptorStore
	options Options
}

func (s *stencilClient) Parse(className string, data []byte) (protoreflect.ProtoMessage, error) {
	desc, ok := s.store.get(className)
	if !ok {
		return nil, ErrNotFound
	}
	m := dynamicpb.NewMessage(desc).New().Interface()
	err := proto.UnmarshalOptions{Resolver: s.store.extensionResolver}.Unmarshal(data, m)
	return m, err
}

func (s *stencilClient) ParseWithRefresh(className string, data []byte) (protoreflect.ProtoMessage, error) {
	m, err := s.Parse(className, data)
	if err != nil || m.ProtoReflect().GetUnknown() == nil {
		return m, err
	}
	if err = s.Refresh(); err != nil {
		return m, err
	}
	return s.Parse(className, data)
}

func (s *stencilClient) Serialize(className string, data interface{}) (bytes []byte, err error) {
	// message to json
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	// get descriptor
	desc, err := s.GetDescriptor(className)
	if err != nil {
		return
	}

	// construct proto message
	m := dynamicpb.NewMessage(desc).New().Interface()
	if err = protojson.Unmarshal(jsonBytes, m); err != nil {
		return bytes, ErrInvalidDescriptor
	}

	// from proto message to bytes[]
	return proto.Marshal(m)
}

func (s *stencilClient) SerializeWithRefresh(className string, data interface{}) (bytes []byte, err error) {
	bytes, err = s.Serialize(className, data)
	if err == nil || (err != ErrNotFound && err != ErrInvalidDescriptor) {
		return
	}

	if err = s.Refresh(); err != nil {
		return bytes, errors.Wrap(err, "error refreshing descriptor")
	}

	return s.Serialize(className, data)
}

func (s *stencilClient) GetDescriptor(className string) (protoreflect.MessageDescriptor, error) {
	desc, ok := s.store.get(className)
	if !ok {
		return nil, ErrNotFound
	}
	return desc, nil
}

func (s *stencilClient) Close() {
	if s.timer != nil {
		s.timer.Close()
	}
}

func (s *stencilClient) Refresh() error {
	var err error
	for _, url := range s.urls {
		err = multierr.Combine(err, s.store.loadFromURI(url, s.options))
	}
	return err
}

func (s *stencilClient) load() error {
	s.options.setDefaults()
	if s.options.AutoRefresh {
		s.timer = setInterval(s.options.RefreshInterval, func() { s.Refresh() })
	}
	err := s.Refresh()
	return err
}

// NewClient creates stencil client. Downloads proto descriptor file from given url and stores the definitions.
// It will throw error if download fails or downloaded file is not fully contained descriptor file
func NewClient(url string, options Options) (Client, error) {
	s := &stencilClient{store: newStore(), urls: []string{url}, options: options}
	err := s.load()
	return s, err
}

// NewMultiURLClient creates stencil client with multiple urls. Downloads proto descriptor file from given urls and stores the definitions.
// If descriptor files from multiple urls has different schema definitions with same name, last downloaded proto descriptor will override previous entries.
// It will throw error if any of the download fails or any downloaded file is not fully contained descriptor file
func NewMultiURLClient(urls []string, options Options) (Client, error) {
	s := &stencilClient{store: newStore(), urls: urls, options: options}
	err := s.load()
	return s, err
}
