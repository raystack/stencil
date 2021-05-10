# Stencil

Stencil is dynamic schema registry for protobuf. [Protobuf](https://developers.google.com/protocol-buffers) is a Google's language-neutral, platform-neutral, extensible mechanism for serializing structured data. The challenge with using generated source code from protobuf is that for every change in proto definition, it requires to recompile dependent services/packages. This approach works for most applications but it's difficult for general purpose applications that needs to operate on arbitrary protobuf schemas. Stencil enables the developers to specifically tackle this problem.

To work with arbitrary proto schema, application need to load proto schema definition at runtime. Protobuf allows you to define a whole proto file using [google.protobuf.FileDescriptorProto](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L62). A [google.protobuf.FileDescriptorSet](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L57) contains list of FileDescriptorProtos. `protoc` can generate descriptorset file from proto files. Stencil heavily make use of this feature to update proto schemas in runtime.


Stencil has two major components.
- Server
- Client

## Stencil server

Stencil Server is written in Go. It provides REST interface for storing and retrieving protobuf descriptorset file.

### Features

 - stores versioned history of proto descriptor file on specified namespace and name
 - enforce backward compatability check on upload by default
 - ability to skip some of the backward compatability checks while upload
 - ability to download proto descriptor files
 - provides metadata API to retrieve latest version number given a name and namespace
 - ability to download latest proto descriptor file
 - support for multiple backend storage services (Local storage, Google cloud storage, S3, Azure blob storage and in-memory storage)

{% page-ref page="server/overview.md" %}

## Stencil client

Stencil client abstracts handling of descriptorset file on client side. Currently we officially support Stencil client in Java, Go, JS languages.

### Features
 - downloading of descriptorset file from server
 - parse API to deserialize protobuf encoded messages
 - lookup API to find proto descriptors
 - inbuilt strategies to refresh protobuf schema definitions.

{% page-ref page="clients/overview.md" %}
