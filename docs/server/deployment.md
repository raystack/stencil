# Stencil server

This doc describes Deployment instructions for stencil server

## Running the server

Run the following command to run from docker image
```bash
$ docker run -e PORT=8000 -e BUCKETURL=file://root -p 8000:8000 odpf/stencil
```

Run the following commands to compile from source
```bash
$ git clone git@github.com:odpf/stencil.git
$ cd stencil/server
$ go build -o stencil
$ ./stencil # specify envs before executing this command
```

## Configuring the stencil server

To run the stencil server, you will need to add the following environment variables

`BUCKETURL`: is common across different backend stores. Please refer URL structure [here](https://gocloud.dev/concepts/urls/) for configuring different backend stores.

`PORT`: port number default to `8080`

Following table represents required variables to authenticate for different backend stores


| Backend store | URL scheme | ENV variables needed to authenticate     | Description                |
| :-------- | :------- | :---------- | :------------------------- |
| Google cloud storage | `gs://` | `GOOGLE_APPLICATION_CREDENTIALS` | Value should point to service account key file. Refer [here](https://cloud.google.com/storage/docs/reference/libraries#setting_up_authentication) to generate key file |
| Azure cloud storage | `azblob://` | `AZURE_STORAGE_ACCOUNT`, `AZURE_STORAGE_KEY`, `AZURE_STORAGE_SAS_TOKEN` | `AZURE_STORAGE_ACCOUNT` is required, along with one of the other two. refer [here](https://gocloud.dev/howto/blob/#azure) for more details |
| AWS cloud storage | `s3://` | refer [here](https://docs.aws.amazon.com/sdk-for-go/api/aws/session/) for list of envs needed | [reference](https://gocloud.dev/howto/blob/#s3) |
| Local storage | `file://` |none | No extra envs required |
| In memory storage | `mem://` | none | No Extra envs required |
