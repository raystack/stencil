# Stencil


![test workflow](https://github.com/odpf/stencil/actions/workflows/server-test.yaml/badge.svg)
![release workflow](https://github.com/odpf/stencil/actions/workflows/release.yml/badge.svg)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?logo=apache)](LICENSE)
[![Version](https://img.shields.io/github/v/release/odpf/stencil?logo=semantic-release)](Version)


Stencil is dynamic schema registry for protobuf. [Protobuf](https://developers.google.com/protocol-buffers) is a Google's language-neutral, platform-neutral, extensible mechanism for serializing structured data. The challenge with using generated source code from protobuf is that for every change in proto definition, it requires to recompile dependent services/packages. This approach works for most applications but it's difficult for general purpose applications that needs to operate on arbitrary protobuf schemas. Stencil enables the developers to specifically tackle this problem.

To work with arbitrary proto schema, application need to load proto schema definition at runtime. Protobuf allows you to define a whole proto file using [google.protobuf.FileDescriptorProto](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L62). A [google.protobuf.FileDescriptorSet](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L57) contains list of FileDescriptorProtos. `protoc` can generate descriptorset file from proto files. Stencil heavily make use of this feature to update proto schemas in runtime.

<p align="center"><img src="./docs/assets/overview.svg" /></p>

## Key Features
Discover why users choose Stencil as their main schema registry

* **Version history** Stencil stores versioned history of proto descriptor file on specified namespace and name
* **Backward compatibility** enforce backward compatability check on upload by default
* **Flexbility** ability to skip some of the backward compatability checks while upload
* **Descriptor fetch** ability to download proto descriptor files
* **Metadata** provides metadata API to retrieve latest version number given a name and namespace
* **Clients in multiple languages** Stencil provides clients in GO, JAVA, JS languages to interact with Stencil server and deserialize messages using dynamic schema


## Usage

Explore the following resources to get started with Stencil:

* [Documentation](http://odpf.gitbook.io/stencil) provides guidance on using stencil.
* [Server](/server) provides details on getting started with stencil server.
* [Clients](/clients) provides reference to supported stencil clients.

## Clients

 - [Java](clients/java)
 - [Go](clients/go)
 - [Javascript](clients/js)

## Contribute

Development of Stencil happens in the open on GitHub, and we are grateful to the community for contributing bugfixes and improvements. Read below to learn how you can take part in improving stencil.

Read our [contributing guide](docs/contribute/contribution.md) to learn about our development process, how to propose bugfixes and improvements, and how to build and test your changes to stencil.

To help you get your feet wet and get you familiar with our contribution process, we have a list of [good first issues](https://github.com/odpf/stencil/labels/good%20first%20issue) that contain bugs which have a relatively limited scope. This is a great place to get started.

## Credits

This project exists thanks to all the [contributors](https://github.com/odpf/stencil/graphs/contributors).

## License
Stencil is [Apache 2.0](LICENSE) licensed.
