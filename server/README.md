
# Stencil server

Stencil is dynamic protobuf schema registry. It provides REST interface for storing and retrieving protobuf file descriptors.




## Features

 - stores versioned history of proto descriptor file on specified namespace and name
 - enforce backward compatability check on upload by default
 - ability to skip some of the backward compatability checks while upload
 - ability to download proto descriptor files
 - provides metadata API to retrieve latest version number given a name and namespace
 - support for multiple backend storage services (Local storage, Google cloud storage, S3, Azure blob storage and in-memory storage)


  
## Installation 


Run the following commands to run from docker image
```bash
$ docker pull odpf/stencil
```

Run the following commands to compile from source
```bash
$ git clone git@github.com:odpf/stencil.git
$ cd stencil/server
$ go build -o stencil
$ ./stencil # specify envs before executing this command
```

### Configuring environment Variables

To run the stencil server, you will need to add the following environment variables

`BUCKETURL` is common across different backend stores. Please refer URL structure [here](https://gocloud.dev/concepts/urls/) for configuring different backend stores.

`PORT` port number default to `8080`

Following table represents required variables to authenticate for different backend stores


| Backend store | ENV variable     | Description                |
| :-------- | :------- | :------------------------- |
| Google cloud storage | `GOOGLE_APPLICATION_CREDENTIALS` | Value should point to service account key file. Refer [here](https://cloud.google.com/storage/docs/reference/libraries#setting_up_authentication) to generate key file |
| Azure cloud storage | `AZURE_STORAGE_ACCOUNT`, `AZURE_STORAGE_KEY`, `AZURE_STORAGE_SAS_TOKEN` | `AZURE_STORAGE_ACCOUNT` is required, along with one of the other two. refer [here](https://gocloud.dev/howto/blob/#azure) for more details |
| AWS cloud storage | refer [here](https://docs.aws.amazon.com/sdk-for-go/api/aws/session/) for list of envs needed | [reference](https://gocloud.dev/howto/blob/#s3) |
| Local storage | none | No extra env required |


## Quick start API usage examples

The following assumes you have Stencil server up and running at port 8080 and `protoc` is installed.

```bash
$ mkdir example
$ cd example
# create example proto file. You can add as many proto files as you want.
$ echo "syntax=\"proto3\";\npackage stencil;\nmessage One {\n  int32 field_one = 1;\n}" > 1.proto

# create descriptor file
$ protoc --descriptor_set_out=./file.desc --include_imports ./1.proto

# upload descriptor file to server with name as `example` under `quickstart` namespace
$ curl -X POST http://localhost:8080/v1/namespaces/quickstart/descriptors -F "file=@./file.desc" -F "version=0.0.1" -F "name=example" -F "latest=true" -H "Content-Type: multipart/form-data"

# get list of descriptors available in a namespace
$ curl -X GET http://localhost:8080/v1/namespaces/quickstart/descriptors

# get list of versions available for particular descriptor
$ curl -X GET http://localhost:8080/v1/namespaces/quickstart/descriptors/example/versions

# download specific version of particular desciptor
$ curl -X GET http://localhost:8080/v1/namespaces/quickstart/descriptors/example/versions/0.0.1

# download latest version of particular descriptor
$ curl -X GET http://localhost:8080/v1/namespaces/quickstart/descriptors/example/versions/latest

# get latest version number of particular descriptor
$ curl -X GET http://localhost:8080/v1/namespaces/quickstart/metadata/example

# modify latest version number of particular descriptor
$ curl -X POST 'http://localhost:8080/v1/namespaces/quickstart' -H 'Content-Type: application/json' --data-raw '{"name": "example","version": "0.0.1"}'
```

