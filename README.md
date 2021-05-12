# Stencil

![test workflow](https://github.com/odpf/stencil/actions/workflows/server-test.yaml/badge.svg)
![release workflow](https://github.com/odpf/stencil/actions/workflows/release.yml/badge.svg)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?logo=apache)](LICENSE)
[![Version](https://img.shields.io/github/v/release/odpf/stencil?logo=semantic-release)](Version)

Stencil is dynamic schema registry for protobuf. Protobuf is a great efficient and fast mechanism for serializing structured data. The challenge with protobuf is that for every change it requires to recompile the package to generate the necessary classes. This is not a big challenge if you have protobuf enclosed in your application and compile at startup. But if you have thousands of protos stored in central registry and 100s of applications use them. Updating dependencies of compiled proto jar can soon become a nightmare.

Protobuf allows you to define a whole proto file using [google.protobuf.FileDescriptorProto](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L62). A [google.protobuf.FileDescriptorSet](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L57) contains list of FileDescriptorProto. Stencil heavily make use of this feature to update proto schemas in runtime.

<p align="center"><img src="./docs/assets/overview.svg" /></p>

## Key Features
Discover why users choose Firehose as their main Kafka Consumer

* **Version history** Stencil stores versioned history of proto descriptor file on specified namespace and name
* **Backward compatibility** enforce backward compatability check on upload by default
* **Flexbility** ability to skip some of the backward compatability checks while upload
* **Descriptor fetch** ability to download proto descriptor files
* **Metadata** provides metadata API to retrieve latest version number given a name and namespace
* **Multiple backends** support for multiple backend storage services (Local storage, Google cloud storage, S3, Azure blob storage and in-memory storage)

## Usage

Explore the following resources to get started with Stencil:

* [Documentation](http://odpf.gitbook.io/stencil) provides guidance on using stencil.
* [Server](/server) provides details on getting started with stencil server.
* [Clients](/clients) provides reference to supported stencil clients.

## Clients

 - [Java](clients/java)
 - [Go](clients/go)
 - [Javascript](clients/js)

### Managing descriptors behind an HTTP Server

For serving the protobuf descriptor set artifacts and their versions use a Stencil Server.
This also helps to easily update the descriptor sets by allowing us to push Protobuf Descriptor sets directly.

#### Endpoints

```http
GET https://stencil-hostname.example.com/artifactory/proto-descriptors/:stencil_repo/:version
PUT https://stencil-hostname.example.com/artifactory/proto-descriptors/:stencil_repo/:version
GET https://stencil-hostname.example.com/metadata/proto-descriptors/:stencil_repo/version
PUT https://stencil-hostname.example.com/metadata/proto-descriptors/:stencil_repo/version
```
The section avoids handling/recommending AUTH mechanism as it is left to the user. The samples below assume that basic Auth set up.

Attribute reference -
 - stencil_repo - A set of protbuf definitions that will be bundled together in a descriptor set
 - version - Artifact version promoted for production use

#### Example Usage

To push a new version of proto descriptor -

```sh
curl -u test:pasword -X PUT "https://stencil-hostname.example.com/artifactory/proto-descriptors/test-stencil-repo/0.0.5" -T /path/to/protobuf/descriptor/set/file
```


To set version of latest promoted artifact -
```sh
curl -u test:password -X PUT "https://stencil-hostname.example.com/metadata/proto-descriptors/test-stencil-repo/version" -d value="0.0.5"
```

#### CI/CD

Use the script `stencil-push.sh` in CI for uploading Protobuf descriptor sets -

```sh
# Should be set these as secrets in CI
export STENCIL_HOSTNAME=stencil-hostname.example.com
export STENCIL_USERNAME=test-stencil-repo
export STENCIL_PASSWORD=test-password

./stencil-push.sh 0.0.1 test-stencil-repo /path/to/protobuf/descriptor/set/file
```

## Contribute

Development of Stencil happens in the open on GitHub, and we are grateful to the community for contributing bugfixes and improvements. Read below to learn how you can take part in improving stencil.

Read our [contributing guide](docs/contribute/contribution.md) to learn about our development process, how to propose bugfixes and improvements, and how to build and test your changes to stencil.

To help you get your feet wet and get you familiar with our contribution process, we have a list of [good first issues](https://github.com/odpf/stencil/labels/good%20first%20issue) that contain bugs which have a relatively limited scope. This is a great place to get started.

## Credits

This project exists thanks to all the [contributors](https://github.com/odpf/stencil/graphs/contributors).

## License
Stencil is [Apache 2.0](LICENSE) licensed.