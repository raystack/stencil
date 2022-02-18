# Introduction

Stencil is a schema registry that provides schema mangement and validation to ensure data compatibility across applications. It enables developers to create, manage and consume schemas dynamically, efficiently, and reliably, and provides a simple way to validate data against those schemas. Stencil support multiple formats including Protobuf, Avro and JSON.

![](/assets/overview.svg)

## What is schema registry?

Schema registry is the sole source of truth for schema definitions. When producers or developers create a specific data schema for use within the system, the data schema is stored in the registry that is accessible to all consumers and relevant parties. A central schema registry

- Provides a serving layer for your metadata to build unified documentation for all schemas.
- Provides an API for storing and retrieving schemas.
- Stores versioned history of all schemas.
- Allows evolution of schema as per compatibility settings and prevents breaking changes by validating the compatibility of schemas.
- Provides clients in multiple languages for serializing and deserializing messages using schemas stored in registry.
- Publishers and consumers use these clients to serialize, deserialize message when producing them to stream, etc.
- Provides multiple backends as a storage system to schemas.
- Help in tackling organisational challenges by data policy enforcement.

## Need for schemas

Schema is a strongly typed description of your data/record that clearly defines the structure, type, and meaning of data. Schemas serve as contracts to enforce communication protocol between consumers and producers using messaging systems. It can also serve as a contract between services as well as between teams to establish a common contract of trust between two parties.
Lack of schema poses multiple challenges. Data producers can make arbitrary changes to the data which can create issues for downstream services to interpret the data. For example:

- The field you're looking for doesn't exist anymore because data producers removed it.
- The type of the field has changed (e.g. what used to be a String is now an Integer)

Schemas in an event-driven architecture should not be an afterthought and should be defined before you start sending the events.

## Data formats

When sending data over the network, you need a way to encode data into bytes. There is a wide variety of data serialization formats, including XML, JSON, BSON, YAML, MessagePack, Protocol Buffers, Thrift, and Avro. 
Multiple programming languages also provide specific serialization methods but using language-specific serialization makes consuming data very hard in different programming languages.

Choice of format for an application is subject to a variety of factors, including data complexity, necessity for humans to read it, latency, and storage space concerns. A good data format should fulfill the following needs.

- The data format should be language and platform-neutral which can be used in communication and data storage.
- The data format should be efficient, provide ease of development, and mature libraries across different programming languages.
- Serialized data should not take a lot of space, so binary formats are a natural fit.
- The data format should have very little protocol parsing time.
- The data format should support the versioning and schema evolution.

## Schema evolution

Schema can be changed, producers and consumers can hold different versions of the scheme simultaneously, and it all continues to work. The schema can be changed, producers and consumers can hold different versions of the scheme simultaneously, and it all continues to work. When working with a large production system, this is a highly valuable feature as it allows you to independently update different elements of the system at different times without concerns about compatibility.
