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
  implementation group: 'org.raystack', name: 'stencil', version: '0.4.1'
```

#### Maven

```xml
<dependency>
  <groupId>org.raystack</groupId>
  <artifactId>stencil</artifactId>
  <version>0.1.6</version>
</dependency>
```

### Creating a stencil Client instance

Stencil client can be created in different modes.

#### Loading Descriptor from Protobuf Class available in the classpath

```java
import org.raystack.stencil.client.StencilClient;
import org.raystack.stencil.StencilClientFactory;

StencilClient stencilClient = StencilClientFactory.getClient();
```

#### Create client with remote URL

```java
import org.raystack.stencil.config.StencilConfig;

String url = "http://url/to/proto/descriptor-set/file";
StencilClient stencilClient = StencilClientFacorty.getClient(url, StencilConfig.builder().build());
```

#### Creating MultiURL client

```java
import org.raystack.stencil.config.StencilConfig;

ArrayList<String> urls = new ArrayList<String>();
urls.add("http://localhost:8082/v1beta1/...");
StencilClient stencilClient = StencilClientFacorty.getClient(urls, StencilConfig.builder().build());
```

#### With StatsD client for monitoring

```java
// From https://github.com/tim-group/java-statsd-client
import com.timgroup.statsd.StatsDClient;
import com.timgroup.statsd.NonBlockingStatsDClient;

StatsDClient statDClient = new NonBlockingStatsDClient("my.prefix", "statsd-host", 8125);
StencilClient stencilClient = StencilClientFactory.getClient(url, StencilConfig.builder().statsDClient(statsDClient).build());
```

#### With Schema Update Listener

Whenever schema has changed this listener will be called.

```java
import org.raystack.stencil.SchemaUpdateListener;

SchemaUpdateListener updateListener = new SchemaUpdateListenerImpl();
StencilClient stencilClient = StencilClientFactory.getClient(url, StencilConfig.builder().updateListener(updateListener).build());
```

#### With version based refresh strategy

If url belongs to stencil server, client can choose to refresh schema data only if there is a new version available.

```java
import org.raystack.stencil.cache.SchemaRefreshStrategy;

StencilConfig config = StencilConfig.builder().refreshStrategy(SchemaRefreshStrategy.versionBasedRefresh()).build();
StencilClient stencilClient = StencilClientFactory.getClient(url, config);
```

#### Passing custom headers

While sending request to specified URL, client can be configured to pass headers as well.

```java
import org.apache.http.Header;
import org.apache.http.HttpHeaders;
import org.apache.http.message.BasicHeader;

Header authHeader = new BasicHeader(HttpHeaders.AUTHORIZATION, "Bearer " + token);
List<Header> headers = new ArrayList<Header>();
headers.add(authHeader);
StencilConfig config = StencilConfig.builder().fetchHeaders(headers).build();
StencilClient stencilClient = StencilClientFactory.getClient(url, config);
```

### Getting descriptor

Given the name of the Proto-Class StencilClient returns the Descriptor for it.

```java
import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;

Descriptors.Descriptor descriptor = stencilClient.get(protoClassName);
```

### Parsing message

```java
import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;

Descriptors.Descriptor descriptor = stencilClient.get(protoClassName);
DynamicMessage message = DynamicMessage.parseFrom(descriptor, bytes);
```

#### Using Parser interface

```java
import org.raystack.stencil.Parser;
import com.google.protobuf.DynamicMessage;

Parser protoParser = stencilClient.getParser("com.example.proto.schema");
DynamicMessage message = protoParser.parse(bytes)
```

### Publishing

The client is published and released via github workflow and uses github tag for versioning.

### Notes

- Stencil uses `java-statsd-client` from `com.timgroup`, Please use the same client in your application for statsd
