# Stencil server

This doc describes Deployment instructions for stencil server

## Running the server

Run the following command to run from docker image
```bash
$ docker run -e PORT=8000 -e DB_CONNECTIONSTRING="postgres://postgres@localhost:5432/db" -p 8000:8000 odpf/stencil
```

Run the following commands to compile from source
```bash
$ git clone git@github.com:odpf/stencil.git
$ cd stencil/server
$ go build -o stencil
$ ./stencil # specify envs before executing this command
```

## Configuring the stencil server

### Configuring environment Variables

To run the stencil server, you will need to add the following environment variables

| ENV          | Description          |
| :------------ | :--------------------- |
| `PORT` | port number default to `8080` |
| `TIMEOUT` | graceful time to wait before shutting down the server. Takes `time.Duration` format. Eg: `30s` or `20m` |
| `DB_CONNECTIONSTRING` | postgres db connection [url](https://www.postgresql.org/docs/11/libpq-connect.html#LIBPQ-CONNSTRING). Eg: `postgres://postgres@localhost:5432/db_name` |
| `NEWRELIC_ENABLED` | boolean to enable newrelic |
| `NEWRELIC_APPNAME` | appname |
| `NEWRELIC_LICENSE` | License key for newrelic |
