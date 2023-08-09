import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";

# Quickstart

This quick start will explore how to use Stencil command line interface and client libraries inside your application code. As part of this quick start we will start stencil server, create schema and then use stencil clients to serialise and deserialise data using registered schemas.

## Prerequisites

- [Docker](../installation#using-docker-image) or a [local installation](../installation#binary-cross-platform) of the Stencil binary.
- A development environment applicable to one of the languages in this quick start (currently Go, Java, and JavaScript).
- Postgres database and [protoc](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation) if your schema format is protobuf.

## Step 1: Start server

<Tabs>
<TabItem value="executable" label="Executable" default>

Run stencil server locally using installed stencil binary.
Note: Below command assumes `stencil_dev` db present in your postgres instance.

```bash
$ export DB_CONNECTIONSTRING=postgres://postgres@localhost:5432/stencil_dev?sslmode=disable

# Run database migrations
$ stencil server migrate

# Stencil server
$ stencil server start

# Check if server running
$ curl -X GET http://localhost:8080/ping
```

</TabItem>
<TabItem value="docker" label="Docker">

Run stencil server locally using docker
Note: Below command assumes `stencil_dev` db present in your postgres instance.

```bash
# Run database migrations
$ docker run -e PORT=8000 -e DB_CONNECTIONSTRING=postgres://postgres@host.docker.internal:5432/stencil_dev?sslmode=disable -p 8000:8000 raystack/stencil server migrate

# Stencil server at port 8000
$ docker run -e PORT=8000 -e DB_CONNECTIONSTRING=postgres://postgres@host.docker.internal:5432/stencil_dev?sslmode=disable -p 8000:8000 raystack/stencil server start

# Check if server running
$ curl -X GET http://localhost:8000/ping
```

</TabItem>
</Tabs>

## Step 2: Create schema

<Tabs>
<TabItem value="protobuf" label="Protobuf">

```bash
$ mkdir example
$ cd example

# Create a sample proto schema.
$ echo "syntax=\"proto3\";
  package stencil;
  message One {
    int32 field_one = 1;
  }" > schema.proto

# Create proto descriptor file
$ protoc --descriptor_set_out=./schema.desc --include_imports ./**/*.proto
```

</TabItem>
<TabItem value="avro" label="Avro">

```bash
$ mkdir example
$ cd example

# Create a sample avro schema.
$ echo "{
   \"type\" : \"record\",
   \"namespace\" : \"Tutorialspoint\",
   \"name\" : \"Employee\",
   \"fields\" : [
      { \"name\" : \"Name\" , \"type\" : \"string\" },
      { \"name\" : \"Age\" , \"type\" : \"int\" }
   ]
}" > schema.json

```

</TabItem>
<TabItem value="json" label="JSON">

```bash
$ mkdir example
$ cd example

# Create example JSON schema file.
$ echo "{
  \"type\":\"object\",
  \"properties\":{
    \"f1\":{
      \"type\":\"string\"
      }
    },
  \"additionalProperties\": false
}" > schema.json

```

</TabItem>
</Tabs>

## Step 3: Upload to server

<Tabs>
<TabItem value="cli" label="CLI">

```bash
# --host does not contain the protocol scheme http:// since they internally use GRPC.

# Create namespace named "quickstart" with backward compatibility enabled
$ stencil namespace create quickstart -c COMPATIBILITY_BACKWARD -f FORMAT_PROTOBUF -d "For quickstart guide" --host localhost:8000

# List namespaces
$ stencil namespace list --host localhost:8000

# Upload generated schema proto descriptor file to server with schema name as `example` under `quickstart` namespace.
$ stencil schema create example --namespace=quickstart –-filePath=schema.desc

# Get list of schemas available in a namespace
$ stencil schema list --host localhost:8000

# Get list of versions available for particular schema. These versions are auto generated.
# Version numbers managed by stencil.
$ stencil schema version example -n quickstart  --host localhost:8000

# Download specific version of particular schema
$ stencil schema get example --version 1 --host localhost:8000

# Download latest version of particular schema
$ stencil schema get example  -n quickstart --host localhost:8000
```

</TabItem>
<TabItem value="api" label="API">

```bash
# Create namespace named "quickstart" with backward compatibility enabled
$ curl -X POST http://localhost:8000/v1beta1/namespaces -H 'Content-Type: application/json' -d '{"id": "quickstart", "format": "FORMAT_PROTOBUF", "compatibility": "COMPATIBILITY_BACKWARD", "description": "For quickstart guide"}'

# List namespaces
$ curl http://localhost:8000/v1beta1/namespaces

# Upload generated schema proto descriptor file to server with schema name as `example` under `quickstart` namespace.
$ curl -X POST http://localhost:8000/v1beta1/namespaces/quickstart/schemas/example --data-binary "@schema.desc"

# Get list of schemas available in a namespace
$ curl -X GET http://localhost:8000/v1beta1/namespaces/quickstart/schemas

# Get list of versions available for particular schema. These versions are auto generated.
# Version numbers managed by stencil.
$ curl -X GET http://localhost:8000/v1beta1/namespaces/quickstart/schemas/example/versions

# Download specific version of particular schema
$ curl -X GET http://localhost:8000/v1beta1/namespaces/quickstart/schemas/example/versions/1

# Download latest version of particular schema
$ curl -X GET http://localhost:8000/v1beta1/namespaces/quickstart/schemas/example;
```

</TabItem>
</Tabs>

## Step 4: Using client

Let's use this API in our GO client

```go
package main

import (
       "log"
       stencil "github.com/raystack/stencil/clients/go"
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
