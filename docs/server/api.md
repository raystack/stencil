


# Stencil server
  

## Informations

### Version

0.1.0

## Tags

  ### <span id="tag-descriptors"></span>descriptors

Manage descriptors

  ### <span id="tag-metadata"></span>metadata

manage latest versions for uploaded descriptor files

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json
  * multipart/form-data

### Produces
  * application/octet-stream
  * application/json

## All endpoints

###  descriptors

  

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /v1/namespaces/{namespace}/descriptors | [get v1 namespaces namespace descriptors](#get-v1-namespaces-namespace-descriptors) | list all available descriptor names under one namespace |
| GET | /v1/namespaces/{namespace}/descriptors/{name}/versions | [get v1 namespaces namespace descriptors name versions](#get-v1-namespaces-namespace-descriptors-name-versions) | list all available versions for specified descriptor |
| GET | /v1/namespaces/{namespace}/descriptors/{name}/versions/{version} | [get v1 namespaces namespace descriptors name versions version](#get-v1-namespaces-namespace-descriptors-name-versions-version) | download specified descriptor file |
| POST | /v1/namespaces/{namespace}/descriptors | [post v1 namespaces namespace descriptors](#post-v1-namespaces-namespace-descriptors) | upload descriptors |
  


###  metadata

  

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /v1/namespaces/{namespace}/metadata/{name} | [get v1 namespaces namespace metadata name](#get-v1-namespaces-namespace-metadata-name) | get latest version for specified descriptor |
| POST | /v1/namespaces/{namespace}/metadata | [post v1 namespaces namespace metadata](#post-v1-namespaces-namespace-metadata) | update metadata |
  


###  operations

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /ping | [ping](#ping) | service health check |
  


## Paths

### <span id="get-v1-namespaces-namespace-descriptors"></span> list all available descriptor names under one namespace (*GetV1NamespacesNamespaceDescriptors*)

```
GET /v1/namespaces/{namespace}/descriptors
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-v1-namespaces-namespace-descriptors-200) | OK | returns list of descriptor names |  | [schema](#get-v1-namespaces-namespace-descriptors-200-schema) |

#### Responses


##### <span id="get-v1-namespaces-namespace-descriptors-200"></span> 200 - returns list of descriptor names
Status: OK

###### <span id="get-v1-namespaces-namespace-descriptors-200-schema"></span> Schema
   
  

[]string

### <span id="get-v1-namespaces-namespace-descriptors-name-versions"></span> list all available versions for specified descriptor (*GetV1NamespacesNamespaceDescriptorsNameVersions*)

```
GET /v1/namespaces/{namespace}/descriptors/{name}/versions
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| name | `path` | string | `string` |  | ✓ |  |  |
| namespace | `path` | string | `string` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-v1-namespaces-namespace-descriptors-name-versions-200) | OK | returns list of versions |  | [schema](#get-v1-namespaces-namespace-descriptors-name-versions-200-schema) |

#### Responses


##### <span id="get-v1-namespaces-namespace-descriptors-name-versions-200"></span> 200 - returns list of versions
Status: OK

###### <span id="get-v1-namespaces-namespace-descriptors-name-versions-200-schema"></span> Schema
   
  

[]string

### <span id="get-v1-namespaces-namespace-descriptors-name-versions-version"></span> download specified descriptor file (*GetV1NamespacesNamespaceDescriptorsNameVersionsVersion*)

```
GET /v1/namespaces/{namespace}/descriptors/{name}/versions/{version}
```

#### Produces
  * application/octet-stream

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| name | `path` | string | `string` |  | ✓ |  |  |
| namespace | `path` | string | `string` |  | ✓ |  |  |
| version | `path` | string | `string` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-v1-namespaces-namespace-descriptors-name-versions-version-200) | OK | download response |  | [schema](#get-v1-namespaces-namespace-descriptors-name-versions-version-200-schema) |

#### Responses


##### <span id="get-v1-namespaces-namespace-descriptors-name-versions-version-200"></span> 200 - download response
Status: OK

###### <span id="get-v1-namespaces-namespace-descriptors-name-versions-version-200-schema"></span> Schema

### <span id="get-v1-namespaces-namespace-metadata-name"></span> get latest version for specified descriptor (*GetV1NamespacesNamespaceMetadataName*)

```
GET /v1/namespaces/{namespace}/metadata/{name}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| name | `path` | string | `string` |  | ✓ |  |  |
| namespace | `path` | string | `string` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-v1-namespaces-namespace-metadata-name-200) | OK | Success response |  | [schema](#get-v1-namespaces-namespace-metadata-name-200-schema) |

#### Responses


##### <span id="get-v1-namespaces-namespace-metadata-name-200"></span> 200 - Success response
Status: OK

###### <span id="get-v1-namespaces-namespace-metadata-name-200-schema"></span> Schema
   
  

[MetadataResponse](#metadata-response)

### <span id="post-v1-namespaces-namespace-descriptors"></span> upload descriptors (*PostV1NamespacesNamespaceDescriptors*)

```
POST /v1/namespaces/{namespace}/descriptors
```

#### Consumes
  * multipart/form-data

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  |  |
| file | `formData` | file | `io.ReadCloser` |  | ✓ |  | descriptorset file to upload |
| latest | `formData` | boolean | `bool` |  |  |  | mark this descriptor file as latest |
| name | `formData` | string | `string` |  | ✓ |  |  |
| skiprules | `formData` | []string | `[]string` |  |  |  | list of rules to skip |
| version | `formData` | string | `string` |  | ✓ |  | version number for descriptor file. This should follow semantic version compatability |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-v1-namespaces-namespace-descriptors-200) | OK | Success response |  | [schema](#post-v1-namespaces-namespace-descriptors-200-schema) |
| [409](#post-v1-namespaces-namespace-descriptors-409) | Conflict | Conflict |  | [schema](#post-v1-namespaces-namespace-descriptors-409-schema) |

#### Responses


##### <span id="post-v1-namespaces-namespace-descriptors-200"></span> 200 - Success response
Status: OK

###### <span id="post-v1-namespaces-namespace-descriptors-200-schema"></span> Schema

##### <span id="post-v1-namespaces-namespace-descriptors-409"></span> 409 - Conflict
Status: Conflict

###### <span id="post-v1-namespaces-namespace-descriptors-409-schema"></span> Schema

### <span id="post-v1-namespaces-namespace-metadata"></span> update metadata (*PostV1NamespacesNamespaceMetadata*)

```
POST /v1/namespaces/{namespace}/metadata
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  |  |
| body | `body` | [MetadataPayload](#metadata-payload) | `models.MetadataPayload` | | ✓ | | specify name and version in payload |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-v1-namespaces-namespace-metadata-200) | OK | Success response |  | [schema](#post-v1-namespaces-namespace-metadata-200-schema) |

#### Responses


##### <span id="post-v1-namespaces-namespace-metadata-200"></span> 200 - Success response
Status: OK

###### <span id="post-v1-namespaces-namespace-metadata-200-schema"></span> Schema

### <span id="ping"></span> service health check (*ping*)

```
GET /ping
```

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#ping-200) | OK | returns pong message |  | [schema](#ping-200-schema) |

#### Responses


##### <span id="ping-200"></span> 200 - returns pong message
Status: OK

###### <span id="ping-200-schema"></span> Schema

## Models

### <span id="metadata-payload"></span> MetadataPayload


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| name | string| `string` |  | |  |  |
| version | string| `string` |  | |  |  |



### <span id="metadata-response"></span> MetadataResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| updated | string| `string` |  | |  |  |
| version | string| `string` |  | |  |  |


