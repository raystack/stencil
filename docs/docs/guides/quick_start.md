# Quick start guide

In this guide we're going to create a example proto file descriptorset then setup local stencil server with local file system as backend, then try out server APIs then we move on to using these APIs in Stencil GO client.

This guide assumes you already have [docker](https://www.docker.com/) and [protoc](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation) installed on your system

## Create proto descriptorset file

```bash
$ mkdir example
$ cd example
# create example proto file. You can add as many proto files as you want.
$ echo "syntax=\"proto3\";\npackage stencil;\nmessage One {\n  int32 field_one = 1;\n}" > one.proto

# create descriptor file
$ protoc --descriptor_set_out=./file.desc --include_imports ./**/*.proto
```

## Run stencil server locally using docker with local filesystem as backend store

```bash
# This will run stencil server at port 8000
$ docker run -e PORT=8000 -e BUCKETURL=file://root -p 8000:8000 odpf/stencil

# check if server running
$ curl -X GET http://localhost:8000/ping
```

## Try Stencil server APIs

```bash
# upload descriptor file to server with name as `example` under `quickstart` namespace
$ curl -X POST http://localhost:8000/v1/namespaces/quickstart/descriptors -F "file=@./file.desc" -F "version=0.0.1" -F "name=example" -F "latest=true" -H "Content-Type: multipart/form-data"

# get list of descriptors available in a namespace
$ curl -X GET http://localhost:8000/v1/namespaces/quickstart/descriptors

# get list of versions available for particular descriptor
$ curl -X GET http://localhost:8000/v1/namespaces/quickstart/descriptors/example/versions

# download specific version of particular desciptor
$ curl -X GET http://localhost:8000/v1/namespaces/quickstart/descriptors/example/versions/0.0.1

# download latest version of particular descriptor
$ curl -X GET http://localhost:8000/v1/namespaces/quickstart/descriptors/example/versions/latest

# get latest version number of particular descriptor
$ curl -X GET http://localhost:8000/v1/namespaces/quickstart/metadata/example

# now let's try uploading breaking proto definition. Note that proto field number has changed from 1 to 2.
$ echo "syntax=\"proto3\";\npackage stencil;\nmessage One {\n  int32 field_one = 2;\n}" > one.proto

# create descriptor file
$ protoc --descriptor_set_out=./file.desc --include_imports ./**/*.proto

# now try to upload this descriptor file with same name as before but different version. This call should fail.
$ curl -X POST http://localhost:8000/v1/namespaces/quickstart/descriptors -F "file=@./file.desc" -F "version=0.0.2" -F "name=example" -F "latest=true" -H "Content-Type: multipart/form-data"

# now let's try fixing our proto add a new field without having any breaking changes.
$ echo "syntax=\"proto3\";\npackage stencil;\nmessage One {\n  int32 field_one = 1;\nint32 field_two = 2;\n}" > one.proto

# create descriptor file
$ protoc --descriptor_set_out=./file.desc --include_imports ./**/*.proto

# now try to upload this descriptor file with same name as before but different version. This call should succeed. Note latest form field as false. We update latest tag using metadata API in next step
$ curl -X POST http://localhost:8000/v1/namespaces/quickstart/descriptors -F "file=@./file.desc" -F "version=0.0.2" -F "name=example" -F "latest=false" -H "Content-Type: multipart/form-data"

# modify latest version number of particular descriptor
$ curl -X POST 'http://localhost:8000/v1/namespaces/quickstart' -H 'Content-Type: application/json' --data-raw '{"name": "example","version": "0.0.2"}'
```

## Let's use this API in our GO client

```go
package main

import (
       "log"
       stencil "github.com/odpf/stencil/clients/go"
)

func main() {
    url := "http://localhost:8000/v1/namespaces/quickstart/descriptors/example/versions/latest"
    client, err := stencil.NewClient(url, stencil.Options{})
    if err != nil {
      log.Fatal("Unable to create client", err)
      return
    }
    desc, err := client.GetDescriptor("stencil.One")
    // ...
}
```
