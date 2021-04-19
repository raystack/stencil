# Java Client for Stencil

Stencil is dynamic schema registry for protobuf. Protobuf is a great efficient and fast mechanism for serializing structured data. The challenge with protobuf is that for every change it requires to recompile the package to generate the necessary classes. This is not a big challenge if you have protobuf enclosed in your application and compile at startup. But if you have thousands of protos stored in central registry and 100s of applications use them. Updating depndencies of compiled proto jar can soon become a nightmare.

Protobuf allows you to define a protobuf file using DescriptorSet. A FileDescriptorSet is basically a description of the proto file i.e. it’s name, it’s package name, it’s dependencies and the messages it contains. Once you have the descriptor file, you can simply read it in any language to create a FileDescriptor Object. Now any serialized ProtoMessage can be deserialized using DynamicMessage and ProtoMessage descriptor.

## Requirements

  - [Gradle v6+](https://gradle.org/)
  - [JDK 8+](https://openjdk.java.net/projects/jdk8/)

## Usage

### Add stencil as dependency

#### Gradle

```groovy
  implementation group: 'io.odpf', name: 'stencil', version: '0.1.0'
```

#### Maven

```xml
<dependency>
  <groupId>io.odpf</groupId>
  <artifactId>stencil</artifactId>
  <version>0.1.0</version>
</dependency>
```

### Creating a stencil Client instance

Stencil client scan be created in different modes.

#### Loading Descriptor from Protobuf Class available in the classpath

```java
import io.odpf.stencil.client.StencilClient;
import io.odpf.stencil.StencilClientFactory;

StencilClient stencilClient = StencilClientFactory.getClient();
```


#### Fetching DescriptorSet file from remote URL

```java
import io.odpf.stencil.config.StencilConfig;

String url = "http://url/to/proto/descriptor-set/file";
StencilClient stencilClient = StencilClientFacorty.getClient(url, StencilConfig.builder().build());
```


#### With StatsD client for monitoring

```java
// From https://github.com/tim-group/java-statsd-client
import com.timgroup.statsd.StatsDClient;
import com.timgroup.statsd.NonBlockingStatsDClient;

StatsDClient statDClient = new NonBlockingStatsDClient("my.prefix", "statsd-host", 8125);
StencilClient stencilClient = StencilClientFactory.getClient(url, StencilConfig.builder().build(), statDClient)
```

### Getting descriptor

Given the name of the Proto-Class StencilClient returns the Descriptor for it.

```java
import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;

Descriptors.Descriptor descriptor = stencilClient.get(protoClassName);

DynamicMessage message = DynamicMessage.parseFrom(descriptor, bytes);
```

### Parsing message

The descriptor obtained above can be used for parsing serialized protobuf message bytes as shown above or bytes can be directly parsed using the `ProtoParser` class, which provides auto update of descriptors in case stencil cache gets updated, as shown below -

```java
import io.odpf.stencil.parser.ProtoParser;
import com.google.protobuf.DynamicMessage;

ProtoParser protoParser = new ProtoParser(stencilClient, "com.example.proto.schema");
DynamicMessage message = protoParser.parse(bytes)
```

### Configurations

```java
// stencil default configs
StencilConfig config = StencilConfig.builder()
        .fetchTimeoutMs(10000)
        .fetchRetries(4)
        .fetchBackoffMinMs(0L)
        // .fetchAuthBearerToken("TOKEN")
        .cacheAutoRefresh(false)
        .cacheTtlMs(0L)
        .build();
```

### Publishing

The client is published and released via github workflow and uses github tag for versioning.

### Notes

- Stencil uses `java-statsd-client` from `com.timgroup`, Please use the same client in your application for statsd
