# Stencil

Stencil is dynamic schema registry for protobuf. Protobuf is a great efficient and fast mechanism for serializing structured data. The challenge with protobuf is that for every change it requires to recompile the package to generate the necessary classes. This is not a big challenge if you have protobuf enclosed in your application and compile at startup. But if you have thousands of protos stored in central registry and 100s of applications use them. Updating dependencies of compiled proto jar can soon become a nightmare.

Protobuf allows you to define a whole proto file using [google.protobuf.FileDescriptorProto](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L62). A [google.protobuf.FileDescriptorSet](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L57) contains list of FileDescriptorProto. Stencil heavily make use of this feature to update proto schemas in runtime.

## Clients

 - [Java](clients/java)
 - [Go](clients/go)
 - TypeScript (planned)

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
