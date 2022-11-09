// Package stencil helps to download and refresh protobuf descriptors from remote server
// and provides helper functions to get protobuf schema descriptors and can parse the messages
// dynamically.
package stencil

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	//ErrNotFound default sentinel error if proto not found
	ErrNotFound = errors.New("not found")
	//ErrInvalidDescriptor is for when descriptor does not match the message
	ErrInvalidDescriptor = errors.New("invalid descriptor")
)

// Client provides utility functions to parse protobuf messages at runtime.
// protobuf messages can be identified by specifying fully qualified generated proto java class name.
type Client interface {
	// Parse parses protobuf message from wire format to protoreflect.ProtoMessage given fully qualified name of proto message.
	// Returns ErrNotFound error if given class name is not found
	Parse(string, []byte) (protoreflect.ProtoMessage, error)
	// Serialize serializes data to bytes given fully qualified name of proto message.
	// Returns ErrNotFound error if given class name is not found
	Serialize(string, interface{}) ([]byte, error)
	// GetDescriptor returns protoreflect.MessageDescriptor given fully qualified proto java class name
	GetDescriptor(string) (protoreflect.MessageDescriptor, error)
	// Close stops background refresh if configured.
	Close()
	// Refresh loads new values from specified url. If the schema is already fetched, the previous value
	// will continue to be used by Parse methods while the new value is loading.
	// If schemas not loaded, then this function will block until the value is loaded.
	Refresh()
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
	// RefreshStrategy refresh strategy to use while fetching schema.
	// Default strategy set to `stencil.LongPollingRefresh` strategy
	RefreshStrategy
	// Logger is the interface used to get logging from stencil internals.
	Logger
}

func (o *Options) setDefaults() {
	if o.RefreshInterval == 0 {
		o.RefreshInterval = 12 * time.Hour
	}
	if o.HTTPOptions.Timeout == 0 {
		o.HTTPOptions.Timeout = 10 * time.Second
	}
}

// NewClient creates stencil client. Downloads proto descriptor file from given url and stores the definitions.
// It will throw error if download fails or downloaded file is not fully contained descriptor file
func NewClient(urls []string, options Options) (Client, error) {
	options.setDefaults()
	stores := []*store{}
	for _, url := range urls {
		s, err := newStore(url, options)
		if err != nil {
			return nil, err
		}
		stores = append(stores, s)
	}

	return &stencilClient{urls: urls, stores: stores, options: options}, nil
}

type stencilClient struct {
	urls    []string
	stores  []*store
	options Options
}

func (s *stencilClient) Parse(className string, data []byte) (protoreflect.ProtoMessage, error) {
	resolver, ok := s.getMatchingResolver(className)
	if !ok {
		return nil, ErrNotFound
	}
	messageType, _ := resolver.Get(className)
	m := messageType.New().Interface()
	err := proto.UnmarshalOptions{Resolver: resolver.GetTypeResolver()}.Unmarshal(data, m)
	return m, err
}

func (s *stencilClient) Serialize(className string, data interface{}) (bytes []byte, err error) {
	// message to json
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	resolver, ok := s.getMatchingResolver(className)
	if !ok {
		return nil, ErrNotFound
	}

	// get descriptor
	messageType, _ := resolver.Get(className)
	// construct proto message
	m := messageType.New().Interface()
	err = protojson.UnmarshalOptions{Resolver: resolver.GetTypeResolver()}.Unmarshal(jsonBytes, m)
	if err != nil {
		return bytes, ErrInvalidDescriptor
	}

	// from proto message to byte[]
	return proto.Marshal(m)
}

func (s *stencilClient) getMatchingResolver(className string) (*Resolver, bool) {
	for _, store := range s.stores {
		resolver, ok := store.getResolver()
		if !ok {
			return nil, false
		}
		_, ok = resolver.Get(className)
		if ok {
			return resolver, ok
		}
	}
	return nil, false
}

func (s *stencilClient) GetDescriptor(className string) (protoreflect.MessageDescriptor, error) {
	resolver, ok := s.getMatchingResolver(className)
	if !ok {
		return nil, ErrNotFound
	}
	desc, _ := resolver.Get(className)
	return desc.Descriptor(), nil
}

func (s *stencilClient) Close() {
	for _, store := range s.stores {
		if store != nil {
			store.Close()
		}
	}
}

func (s *stencilClient) Refresh() {
	var wg sync.WaitGroup
	for _, st := range s.stores {
		wg.Add(1)
		go func(s *store) {
			defer wg.Done()
			s.refresh()
		}(st)
	}
	wg.Wait()
}
