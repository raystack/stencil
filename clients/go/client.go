package stencilclient

import (
	"errors"
	"io"
	"time"

	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

var (
	ErrNotFound = errors.New("not found")
)

// StencilClient provides utility functions to parse protobuf messages at runtime.
// protobuf messages can be identified by specifying fully qualified generated proto java class name.
type StencilClient interface {
	// Parse parses protobuf message from wire format to protoreflect.Message given fully qualified name of proto message.
	// Returns ErrNotFound error if given class name is not found
	Parse(string, []byte) (protoreflect.ProtoMessage, error)
	// GetDescriptor returns protoreflect.MessageDescriptor given fully qualified proto java class name
	GetDescriptor(string) (protoreflect.MessageDescriptor, error)
	// Close stops background refresh if configured.
	Close()
}

// HTTPOptions options for http client
type HTTPOptions struct {
	// Timeout specifies a time limit for requests made by this client
	Timeout time.Duration
	// Headers provide extra headers to be added in requests made by this client
	Headers map[string]string
}

// Options options for stencil client
type Options struct {
	// AutoRefresh boolean to enable or disable autorefresh. Default to false
	AutoRefresh bool
	// RefreshInterval refresh interval to fetch descriptor file from server.
	RefreshInterval time.Duration
	// HTTPOptions options for http client
	HTTPOptions
}

type stencilclient struct {
	timer io.Closer
	urls  []string
	store *descriptorStore
}

func (s *stencilclient) Parse(className string, data []byte) (protoreflect.ProtoMessage, error) {
	desc, ok := s.store.get(className)
	if !ok {
		return nil, ErrNotFound
	}
	m := dynamicpb.NewMessage(desc).New().Interface()
	err := proto.UnmarshalOptions{Resolver: s.store.extensionResolver}.Unmarshal(data, m)
	return m, err
}

func (s *stencilclient) GetDescriptor(className string) (protoreflect.MessageDescriptor, error) {
	desc, ok := s.store.get(className)
	if !ok {
		return nil, ErrNotFound
	}
	return desc, nil
}

func (s *stencilclient) Close() {
	if s.timer != nil {
		s.timer.Close()
	}
}

func (s *stencilclient) refresh(opts Options) error {
	var err error
	for _, url := range s.urls {
		err = multierr.Combine(err, s.store.loadFromURI(url, opts))
	}
	return err
}

func (s *stencilclient) load(opts Options) error {
	if opts.AutoRefresh {
		s.timer = setInterval(opts.RefreshInterval, func() { s.refresh(opts) })
	}
	err := s.refresh(opts)
	return err
}

// NewClient creates stencil client
func NewClient(url string, options Options) (StencilClient, error) {
	s := &stencilclient{store: newStore(), urls: []string{url}}
	err := s.load(options)
	return s, err
}

// NewMultiURLClient creates stencil client with multiple urls
func NewMultiURLClient(urls []string, options Options) (StencilClient, error) {
	s := &stencilclient{store: newStore(), urls: urls}
	err := s.load(options)
	return s, err
}
