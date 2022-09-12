# Overview

Stencil is dynamic protobuf schema registry. It provides REST interface for storing and retrieving protobuf file descriptors.

## Features

- stores versioned history of proto descriptor file on specified namespace and name
- enforce backward compatibility check on upload by default
- ability to skip some of the backward compatibility checks while upload
- ability to download fully contained proto descriptor file for specified proto message [fullName](https://pkg.go.dev/google.golang.org/protobuf@v1.27.1/reflect/protoreflect#FullName)
- provides metadata API to retrieve latest version number given a name and namespace

## Requirements

- postgres 13

## Installation

Run the following commands to run from docker image

```bash
$ docker pull odpf/stencil
```

Run the following commands to compile from source

```bash
$ git clone git@github.com:odpf/stencil.git
$ cd stencil
$ go build -o stencil
$ ./stencil --help

# Create a sample config file.
$ cp config/config.yaml config.yaml

$ ./stencil server start -c config.yaml

```

### Configuring environment Variables

You can also specfify stencil server configurations through following environment variables.
Note: ENV vars takes more precendence over config file.

| ENV                   | Description                                                                                                                                            |
| :-------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------- |
| `PORT`                | port number default to `8080`                                                                                                                          |
| `TIMEOUT`             | graceful time to wait before shutting down the server. Takes `time.Duration` format. Eg: `30s` or `20m`                                                |
| `DB_CONNECTIONSTRING` | postgres db connection [url](https://www.postgresql.org/docs/11/libpq-connect.html#LIBPQ-CONNSTRING). Eg: `postgres://postgres@localhost:5432/db_name` |
| `NEWRELIC_ENABLED`    | boolean to enable newrelic                                                                                                                             |
| `NEWRELIC_APPNAME`    | appname                                                                                                                                                |
| `NEWRELIC_LICENSE`    | License key for newrelic                                                                                                                               |

## Reference

- [API](../reference/api.md)
- [Rules](./rules.md)

## Quick start API usage examples

The following assumes you have Stencil server up and running at port 8080 and `protoc` is installed.

```bash
$ mkdir example
$ cd example
# create example proto file. You can add as many proto files as you want.
$ echo "syntax=\"proto3\";\npackage stencil;\nmessage One {\n  int32 field_one = 1;\n}" > 1.proto

# create descriptor file
$ protoc --descriptor_set_out=./file.desc --include_imports ./1.proto

# create namespace named "quickstart" with backward compatibility enabled
curl -X POST http://localhost:8000/v1beta1/namespaces -H 'Content-Type: application/json' -d '{"id": "quickstart", "format": "FORMAT_PROTOBUF", "compatibility": "COMPATIBILITY_BACKWARD", "description": "This field can be used to store namespace description"}'

# list namespaces
curl http://localhost:8000/v1beta1/namespaces

# upload generated proto descriptor file to server with schema name as `example` under `quickstart` namespace.
curl -X POST http://localhost:8000/v1beta1/namespaces/quickstart/schemas/example --data-binary "@file.desc"

# get list of schemas available in a namespace
curl -X GET http://localhost:8000/v1beta1/namespaces/quickstart/schemas

# get list of versions available for particular schema. These versions are auto generated. Version numbers managed by stencil.
curl -X GET http://localhost:8000/v1beta1/namespaces/quickstart/schemas/example/versions

# download specific version of particular schema
curl -X GET http://localhost:8000/v1beta1/namespaces/quickstart/schemas/example/versions/1

# download latest version of particular schema
curl -X GET http://localhost:8000/v1beta1/namespaces/quickstart/schemas/example;

# now let's try uploading breaking proto definition. Note that proto field number has changed from 1 to 2.
echo "syntax=\"proto3\";\npackage stencil;\nmessage One {\n  int32 field_one = 2;\n}" > one.proto;

# create descriptor file
protoc --descriptor_set_out=./file.desc --include_imports ./**/*.proto;

# now try to upload this descriptor file with same name as before. This call should fail, giving you reason it has failed.
curl -X POST http://localhost:8000/v1/namespaces/quickstart/schemas --data-binary "@file.desc";

# now let's try fixing our proto add a new field without having any breaking changes.
echo "syntax=\"proto3\";\npackage stencil;\nmessage One {\n  int32 field_one = 1;\nint32 field_two = 2;\n}" > one.proto;

# create descriptor file
protoc --descriptor_set_out=./file.desc --include_imports ./**/*.proto

# now try to upload this descriptor file with same name as before. This call should succeed
curl -X POST http://localhost:8000/v1/namespaces/quickstart/schemas --data-binary "@file.desc"

# now try versions api. It should have 2 versions now.
curl -X GET http://localhost:8000/v1beta1/namespaces/quickstart/schemas/example/versions

# upload schema can be called multiple times. Stencil server will retain old version if it's already uploaded. This call won't create new version again. You can verify by using versions API again.
curl -X POST http://localhost:8000/v1/namespaces/quickstart/schemas --data-binary "@file.desc"
```
