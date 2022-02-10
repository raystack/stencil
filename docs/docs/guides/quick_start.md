# Quick start guide

In this guide we're going to create a example proto file descriptorset then setup local stencil server with local file system as backend, then try out server APIs then we move on to using these APIs in Stencil GO client.

This guide assumes you already have [docker](https://www.docker.com/), `postgres` and [protoc](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation) installed on your system.

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

Note: Below command assumes `stencil_dev` db present in your postgres instance.

```bash
# To run migrations
$ docker run -e PORT=8000 -e DB_CONNECTIONSTRING=postgres://postgres@host.docker.internal:5432/stencil_dev?sslmode=disable -p 8000:8000 odpf/stencil server migrate
# This will run stencil server at port 8000
$ docker run -e PORT=8000 -e DB_CONNECTIONSTRING=postgres://postgres@host.docker.internal:5432/stencil_dev?sslmode=disable -p 8000:8000 odpf/stencil server start

# check if server running
$ curl -X GET http://localhost:8000/ping
```

## Try Stencil server APIs

```bash
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

## Let's use this API in our GO client

```go
package main

import (
       "log"
       stencil "github.com/odpf/stencil/clients/go"
)

func main() {
    url := "http://localhost:8000/v1/namespaces/quickstart/descriptors/example/versions/latest"
    client, err := stencil.NewClient([]string{url}, stencil.Options{})
    if err != nil {
      log.Fatal("Unable to create client", err)
      return
    }
    desc, err := client.GetDescriptor("stencil.One")
    // ...
}
```
