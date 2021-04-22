# stencil

```go
import "github.com/odpf/stencil/clients/go"
```

Package stencil helps to download and refresh protobuf descriptors from remote server and provides helper functions to get protobuf schema descriptors and can parse the messages dynamically\.

## Index

- [Variables](<#variables>)
- [type Client](<#type-client>)
  - [func NewClient(url string, options Options) (Client, error)](<#func-newclient>)
  - [func NewMultiURLClient(urls []string, options Options) (Client, error)](<#func-newmultiurlclient>)
- [type HTTPOptions](<#type-httpoptions>)
- [type Options](<#type-options>)


## Variables

```go
var (
    //ErrNotFound default sentinel error if proto not found
    ErrNotFound = errors.New("not found")
)
```

## type [Client](<https://github.com/odpf/stencil/blob/master/clients/go/client.go#L24-L38>)

Client provides utility functions to parse protobuf messages at runtime\. protobuf messages can be identified by specifying fully qualified generated proto java class name\.

```go
type Client interface {
    // Parse parses protobuf message from wire format to protoreflect.ProtoMessage given fully qualified name of proto message.
    // Returns ErrNotFound error if given class name is not found
    Parse(string, []byte) (protoreflect.ProtoMessage, error)
    // ParseWithRefresh parses protobuf message from wire format to `protoreflect.ProtoMessage` given fully qualified name of proto message.
    // Refreshes proto definitions if parsed message has unknown fields and parses the message again.
    // Returns ErrNotFound error if given class name is not found.
    ParseWithRefresh(string, []byte) (protoreflect.ProtoMessage, error)
    // GetDescriptor returns protoreflect.MessageDescriptor given fully qualified proto java class name
    GetDescriptor(string) (protoreflect.MessageDescriptor, error)
    // Close stops background refresh if configured.
    Close()
    // Refresh downloads latest proto definitions
    Refresh() error
}
```

### func [NewClient](<https://github.com/odpf/stencil/blob/master/clients/go/client.go#L129>)

```go
func NewClient(url string, options Options) (Client, error)
```

NewClient creates stencil client

### func [NewMultiURLClient](<https://github.com/odpf/stencil/blob/master/clients/go/client.go#L136>)

```go
func NewMultiURLClient(urls []string, options Options) (Client, error)
```

NewMultiURLClient creates stencil client with multiple urls

## type [HTTPOptions](<https://github.com/odpf/stencil/blob/master/clients/go/client.go#L41-L47>)

HTTPOptions options for http client

```go
type HTTPOptions struct {
    // Timeout specifies a time limit for requests made by this client. Default to 10s.
    // `0` duration not allowed. Client will set to default value (i.e. 10s).
    Timeout time.Duration
    // Headers provide extra headers to be added in requests made by this client
    Headers map[string]string
}
```

## type [Options](<https://github.com/odpf/stencil/blob/master/clients/go/client.go#L50-L58>)

Options options for stencil client

```go
type Options struct {
    // AutoRefresh boolean to enable or disable autorefresh. Default to false
    AutoRefresh bool
    // RefreshInterval refresh interval to fetch descriptor file from server. Default to 12h.
    // `0` duration not allowed. Client will set to default value (i.e. 12h).
    RefreshInterval time.Duration
    // HTTPOptions options for http client
    HTTPOptions
}
```
