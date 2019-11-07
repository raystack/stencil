# Stencil

Dynamic schema registry for protobuf.

## Motivation

Protobuf is a great efficient and fast mechanism for serializing structured data. The challenge with protobuf is that for every change it requires to recompile the package to generate the necessary classes. This is not a big challenge if you have protobuf enclosed in your application and compile at startup. But if you have thousands of protos stored in central registry and 100s of applications use them. Updating depndencies of compiled proto jar can soon become a nightmare.

## How it works

Protobuf allows you to define a protobuf file using DescriptorSet. A FileDescriptorSet is basically a description of the proto file i.e. it’s name, it’s package name, it’s dependencies and the messages it contains. Once you have the descriptor file, you can simply read it in any language to create a FileDescriptor Object. Now any serialized ProtoMessage can be deserialized using DynamicMessage and ProtoMessage descriptor.

## Usage

### Add stencil as gradle dependency

- `compile group: 'com.gojek.de', name: 'stencil-client', version: '2.0.14'`

### Creating a stencil Client instance

Stencil client scan be created in different modes.

- Basic mode

```java
  StencilClient  stencilClient = StencilClientFactory.getClient();
```

This loads the Protobuf Class from the Classpath.

- Descriptor URL path

```java
StecilClient stecnilClient = StencilClientFacorty.getClient(url, Collections.emptyMap());
```

This fetches the artifacts from the provided url.

- With statsd client

```java
 StencilClient stencilClient = StencilClientFactory.getClient(url, new HashMap<>(), stasdClient)
```

### Getting descriptor

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

### Configurations

```
STENCIL_TIMEOUT_MS (10000)
STENCIL_BACKOFF_MS (2000-4000)
STENCIL_RETRIES (4)
TTL_IN_MINUTES (30-60)
```

### Publishing

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

### Notes

- Stencil uses `java-statsd-client` from `com.timgroup`, Please use the same client in your application for statsd
