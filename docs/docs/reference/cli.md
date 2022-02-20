# CLI

## `stencil completion [bash|zsh|fish|powershell]`

Generate shell completion scripts

## `stencil namespace`

Manage namespace

### `stencil namespace create [flags]`

Create a namespace

```
-c, --comp string schema compatibility
-d, --desc string description
-f, --format string schema format
--host string stencil host address eg: localhost:8000
```

### `stencil namespace delete [flags]`

Delete a namespace

```
--host string stencil host address eg: localhost:8000
```

### `stencil namespace get [flags]`

View a namespace

```
--host string stencil host address eg: localhost:8000
```

### `stencil namespace list [flags]`

List all namespaces

```
--host string stencil host address eg: localhost:8000
```

### `stencil namespace update [flags]`

Update a namespace

```
-c, --comp string schema compatibility
-d, --desc string description
-f, --format string schema format
--host string stencil host address eg: localhost:8000
```

## `stencil schema`

Manage schema

### `stencil schema create [flags]`

Create a schema

```
-c, --comp string schema compatibility
-F, --filePath string path to the schema file
-f, --format string schema format
--host string stencil host address eg: localhost:8000
-n, --namespace string parent namespace ID
```

### `stencil schema delete [flags]`

Delete a schema

```
    --host string        stencil host address eg: localhost:8000

-n, --namespace string parent namespace ID
-v, --version int32 particular version to be deleted
```

### `stencil schema diff [flags]`

Diff(s) of two schema versions

```
    --earlier-version int32   earlier version of the schema
    --fullname string         only required for FORMAT_PROTO. fullname of proto schema eg: odpf.common.v1.Version
    --host string             stencil host address eg: localhost:8000
    --later-version int32     later version of the schema

-n, --namespace string parent namespace ID
```

### `stencil schema get [flags]`

View a schema

```
    --host string        stencil host address eg: localhost:8000

-m, --metadata set this flag to get metadata
-n, --namespace string parent namespace ID
-o, --output string path to the output file
-v, --version int32 version of the schema
```

### `stencil schema graph [flags]`

Generate file descriptorset dependencies graph

```
    --host string        stencil host address eg: localhost:8000

-n, --namespace string provide namespace/group or entity name
-o, --output string write to .dot file (default "./proto_vis.dot")
-v, --version int32 provide version number
```

### `stencil schema list [flags]`

List all schemas

```
--host string stencil host address eg: localhost:8000
```

### `stencil schema print [flags]`

Prints snapshot details into .proto files

```
    --filter-path string   filter protocol buffer files by path prefix, e.g., --filter-path=google/protobuf
    --host string          stencil host address eg: localhost:8000

-n, --namespace string provide namespace/group or entity name
-o, --output string the directory path to write the descriptor files, default is to print on stdout
-s, --schema string provide proto repo name
-v, --version int32 provide version number
```

### `stencil schema update [flags]`

Edit a schema

```
-c, --comp string schema compatibility
--host string stencil host address eg: localhost:8000
-n, --namespace string parent namespace ID
```

### `stencil schema version [flags]`

Version(s) of a schema

```
    --host string        stencil host address eg: localhost:8000

-n, --namespace string parent namespace ID
```

## `stencil server <command>`

Server management

### `stencil server migrate [flags]`

Run database migrations

```
-c, --config string Config file path (default "./config.yaml")
```

### `stencil server start [flags]`

Start the server

```
-c, --config string Config file path (default "./config.yaml")
```
