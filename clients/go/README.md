# Stencil go client

Stencil go client package provides a store to lookup protobuf descriptors and options to keep the protobuf descriptors upto date.

It has following features
 - Deserialize protobuf messages directly by specifying protobuf message name
 - Ability to refresh protobuf descriptors in specified intervals
 - Support to download descriptors from multiple urls
## Requirements

 - go 1.16

## Installation

Use `go get`
```
go get github.com/odpf/stencil/clients/go
```

Then import the stencil package into your own code as mentioned below
```go
import stencil "github.com/odpf/stencil/clients/go"
```

## Usage and Documentation
[![Go Reference](https://pkg.go.dev/badge/github.com/odpf/stencil/clients/go.svg)](https://pkg.go.dev/github.com/odpf/stencil/clients/go)

### Creating a client

```go
import stencil "github.com/odpf/stencil/clients/go"

url := "http://url/to/proto/descriptorset/file"
client, err := stencil.NewClient(url, stencil.Options{})
```

### Creating a multiURLClient

```go
import stencil "github.com/odpf/stencil/clients/go"

urls := []string{"http://urlA", "http://urlB"}
client, err := stencil.NewMultiURLClient(urls, stencil.Options{})
```

### Get Descriptor
```go
import stencil "github.com/odpf/stencil/clients/go"

url := "http://url/to/proto/descriptorset/file"
client, err := stencil.NewClient(url, stencil.Options{})
if err != nil {
    return
}
desc, err := client.GetDescriptor("google.protobuf.DescriptorProto")
```

### Parse protobuf message. 
```go
import stencil "github.com/odpf/stencil/clients/go"

url := "http://url/to/proto/descriptorset/file"
client, err := stencil.NewClient(url, stencil.Options{})
if err != nil {
    return
}
data := []byte("")
desc, err := client.Parse("google.protobuf.DescriptorProto", data)
```

Refer to [go documentation](https://pkg.go.dev/github.com/odpf/stencil/clients/go) for all available methods and options.
