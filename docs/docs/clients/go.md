# Go

[![Go Reference](https://pkg.go.dev/badge/github.com/odpf/stencil/clients/go.svg)](https://pkg.go.dev/github.com/odpf/stencil/clients/go)

Stencil go client package provides a store to lookup protobuf descriptors and options to keep the protobuf descriptors upto date.

It has following features

- Deserialize protobuf messages directly by specifying protobuf message name
- Serialize data by specifying protobuf message name
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

## Usage

### Creating a client

```go
import stencil "github.com/odpf/stencil/clients/go"

url := "http://localhost:8000/v1beta1/namespaces/{test-namespace}/schemas/{schema-name}"
client, err := stencil.NewClient([]string{url}, stencil.Options{})
```

### Get Descriptor

```go
import stencil "github.com/odpf/stencil/clients/go"

url := "http://localhost:8000/v1beta1/namespaces/{test-namespace}/schemas/{schema-name}"
client, err := stencil.NewClient([]string{url}, stencil.Options{})
if err != nil {
    return
}
desc, err := client.GetDescriptor("google.protobuf.DescriptorProto")
```

### Parse protobuf message.

```go
import stencil "github.com/odpf/stencil/clients/go"

url := "http://localhost:8000/v1beta1/namespaces/{test-namespace}/schemas/{schema-name}"
client, err := stencil.NewClient([]string{url}, stencil.Options{})
if err != nil {
    return
}
data := []byte("")
parsedMsg, err := client.Parse("google.protobuf.DescriptorProto", data)
```

### Serialize data.

```go
import stencil "github.com/odpf/stencil/clients/go"

url := "http://url/to/proto/descriptorset/file"
client, err := stencil.NewClient([]string{url}, stencil.Options{})
if err != nil {
    return
}
data := map[string]interface{}{}
serializedMsg, err := client.Serialize("google.protobuf.DescriptorProto", data)
```

### Enable auto refresh of schemas

```go
import stencil "github.com/odpf/stencil/clients/go"

url := "http://localhost:8000/v1beta1/namespaces/{test-namespace}/schemas/{schema-name}"
// Configured to refresh schema every 12 hours
client, err := stencil.NewClient([]string{url}, stencil.Options{AutoRefresh: true, RefreshInterval: time.Hours * 12})
if err != nil {
    return
}
desc, err := client.GetDescriptor("google.protobuf.DescriptorProto")
```

### Using VersionBasedRefresh strategy

```go
import stencil "github.com/odpf/stencil/clients/go"

url := "http://localhost:8000/v1beta1/namespaces/{test-namespace}/schemas/{schema-name}"
// Configured to refresh schema every 12 hours
client, err := stencil.NewClient([]string{url}, stencil.Options{AutoRefresh: true, RefreshInterval: time.Hours * 12, RefreshStrategy: stencil.VersionBasedRefresh})
if err != nil {
    return
}
desc, err := client.GetDescriptor("google.protobuf.DescriptorProto")
```

Refer to [go documentation](https://pkg.go.dev/github.com/odpf/stencil/clients/go) for all available methods and options.
