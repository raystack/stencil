# Stencil

Stencil is dynamic schema registry for protobuf. Protobuf is a great efficient and fast mechanism for serializing structured data. The challenge with protobuf is that for every change it requires to recompile the package to generate the necessary classes. This is not a big challenge if you have protobuf enclosed in your application and compile at startup. But if you have thousands of protos stored in central registry and 100s of applications use them. Updating depndencies of compiled proto jar can soon become a nightmare.

Protobuf allows you to define a protobuf file using DescriptorSet. A FileDescriptorSet is basically a description of the proto file i.e. it’s name, it’s package name, it’s dependencies and the messages it contains. Once you have the descriptor file, you can simply read it in any language to create a FileDescriptor Object. Now any serialized ProtoMessage can be deserialized using DynamicMessage and ProtoMessage descriptor.

### Usage

#### Add stencil as gradle dependency

- `compile group: 'com.gojek.de', name: 'stencil', version: '2.0.14'`

#### Creating a stencil Client instance

Stencil client scan be created in different modes.

- Basic mode

```java
  StencilClient  stencilClient = StencilClientFactory.getClient();
```

This loads the Protobuf Class from the Classpath.

- Descriptor URL path

```java
StencilClient stencilClient = StencilClientFacorty.getClient(url, Collections.emptyMap());
```

This fetches the artifacts from the provided url.

- With statsd client

```java
 StencilClient stencilClient = StencilClientFactory.getClient(url, new HashMap<>(), stasdClient)
```

#### Getting descriptor

Given the name of the Proto-Class StencilClient returns the Descriptor for it.

```java
import com.google.protobuf.Descriptors;
public Descriptors.Descriptor descriptor = stencilClient.get(protoClassName);
```

The descriptor obtained above can be used for parsing Consumer Record Value and generating Dynamic Message

```java
DynamicMessage.parseFrom(descriptor, bytes);
OR
ProtoParser protoParser = new ProtoParser(stencilClient, appConfig.getProtoSchema());
protoParser.parse(byte[])
```

#### Configurations

```
STENCIL_TIMEOUT_MS (10000)
STENCIL_BACKOFF_MS (2000-4000)
STENCIL_RETRIES (4)
TTL_IN_MINUTES (30-60)
```

#### Publishing

In order to publish to central maven, you require sonatype credentials and [GnuPG](http://gnupg.org) setup.

To get sonatype credentials, Register to `https://issues.sonatype.org/`  
To setup GnuPG:
Install gpg with `brew install gpg`
Refer to this [url](https://docs.gradle.org/current/userguide/signing_plugin.html#sec:signatory_credentials) to setup GnuPG

After this, add these values to `gradle.properties` in your user directory

```

signing.keyId=<last eight symbols of gnupg keyId>
signing.password=<your passphrase to unlock gpg secrets>
signing.secretKeyRingFile=/Users/me/.gnupg/secring.gpg

ossrhUsername=your-jira-id
ossrhPassword=your-jira-password
```

Upload your gpg keys to ubuntu opengpg server (Required once)
Run the following command:

```
gpg --keyserver hkp://pool.sks-keyservers.net --recv-keys <last eight symbols of gnupg keyId>
gpg --keyserver hkp://keyserver.ubuntu.com --send-keys <last eight symbols of gnupg keyId>
```

### Stencil Server API

For serving the protobuf descriptor set artifacts and their versions we use a Stencil Server.
This also helps us to easily update the descriptor sets by allowing us to push Protobuf Descriptor sets directly.

#### Endpoints

```http
GET https://stencil-hostname.example.com/artifactory/proto-descriptors/:stencil_repo/:version
PUT https://stencil-hostname.example.com/artifactory/proto-descriptors/:stencil_repo/:version
GET https://stencil-hostname.example.com/metadata/proto-descriptors/:stencil_repo/version
PUT https://stencil-hostname.example.com/metadata/proto-descriptors/:stencil_repo/version
```
All above endpoints are behind HTTP Basic Auth.

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

#### Notes

- Stencil uses `java-statsd-client` from `com.timgroup`, Please use the same client in your application for statsd
