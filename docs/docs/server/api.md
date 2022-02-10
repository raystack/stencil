# API
## Version: 0.1.6

### /v1beta1/namespaces

#### GET
##### Summary

List names of namespaces

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1ListNamespacesResponse](#v1beta1listnamespacesresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### POST
##### Summary

Create namespace entry

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body |  | Yes | [v1beta1CreateNamespaceRequest](#v1beta1createnamespacerequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1CreateNamespaceResponse](#v1beta1createnamespaceresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/namespaces/{id}

#### GET
##### Summary

Get namespace by id

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1GetNamespaceResponse](#v1beta1getnamespaceresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### DELETE
##### Summary

Delete namespace by id

##### Description

Ensure all schemas under this namespace is deleted, otherwise it will throw error

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1DeleteNamespaceResponse](#v1beta1deletenamespaceresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### PUT
##### Summary

Update namespace entity by id

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string |
| body | body |  | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1UpdateNamespaceResponse](#v1beta1updatenamespaceresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/namespaces/{id}/schemas

#### GET
##### Summary

List schemas under the namespace

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path |  | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1ListSchemasResponse](#v1beta1listschemasresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/namespaces/{namespaceId}/schemas/{schemaId}
#### POST

##### Summary

Create schema

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| namespaceId | path |  | Yes | string |
| schemaId | path |  | Yes | string |
| X-Format | headers | Can be used to override schema format defined at namespace level | No | [SchemaFormat](#schemaformat) |
| X-Compatibility | headers | Can be used to override schema compatibility defined at namespace level | No | [SchemaCompatibility](#schemacompatibility) |
| body | body | schema data | Yes | valid fileDescriptorSet data/avro/json schemas |

##### Responses
| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1CreateSchemaResponse](#v1beta1CreateSchemaResponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### GET

##### Summary

Get latest schema

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| namespaceId | path |  | Yes | string |
| schemaId | path |  | Yes | string |

##### Responses
| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. Based on schema format, response will return different content types. For avro and json schemas response type is `application/json`. For protobuf response type is `application/octet-stream` | json or byte data |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### DELETE
##### Summary

Delete specified schema

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| namespaceId | path |  | Yes | string |
| schemaId | path |  | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1DeleteSchemaResponse](#v1beta1deleteschemaresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### PATCH
##### Summary

Update only schema metadata

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| namespaceId | path |  | Yes | string |
| schemaId | path |  | Yes | string |
| body | body |  | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1UpdateSchemaMetadataResponse](#v1beta1updateschemametadataresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/namespaces/{namespaceId}/schemas/{schemaId}/meta

#### GET
##### Summary

Create schema under the namespace. Returns version number, unique ID and location

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| namespaceId | path |  | Yes | string |
| schemaId | path |  | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1GetSchemaMetadataResponse](#v1beta1getschemametadataresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/namespaces/{namespaceId}/schemas/{schemaId}/versions

#### GET
##### Summary

List all version numbers for schema

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| namespaceId | path |  | Yes | string |
| schemaId | path |  | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1ListVersionsResponse](#v1beta1listversionsresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/namespaces/{namespaceId}/schemas/{schemaId}/versions/{versionId}

#### GET

##### Summary

Get schema for specified version

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| namespaceId | path |  | Yes | string |
| schemaId | path |  | Yes | string |
| versionId | path |  | Yes | integer |

##### Responses
| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. Based on schema format, response will return different content types. For avro and json schemas response type is `application/json`. For protobuf response type is `application/octet-stream` | json or byte data |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

#### DELETE
##### Summary

Delete specified version of the schema

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| namespaceId | path |  | Yes | string |
| schemaId | path |  | Yes | string |
| versionId | path |  | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1DeleteVersionResponse](#v1beta1deleteversionresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### /v1beta1/search

#### GET
##### Summary

Global Search API

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| namespaceId | query |  | No | string |
| schemaId | query |  | No | string |
| query | query |  | No | string |
| history | query |  | No | boolean |
| versionId | query |  | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1beta1SearchResponse](#v1beta1searchresponse) |
| default | An unexpected error response. | [rpcStatus](#rpcstatus) |

### Models

#### SchemaCompatibility

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| SchemaCompatibility | string |  |  |

Enumarated values for Schema compatibility

| Values |
| ---- |
| COMPATIBILITY_BACKWARD |
| COMPATIBILITY_FORWARD |
| COMPATIBILITY_FULL |

#### SchemaFormat

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| SchemaFormat | string |  |  |
 
Enumarated schema format values

| Values |
| ---- |
| FORMAT_PROTOBUF |
| FORMAT_AVRO |
| FORMAT_JSON |

#### protobufAny

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| typeUrl | string |  | No |
| value | byte |  | No |

#### rpcStatus

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |
| message | string |  | No |
| details | [ [protobufAny](#protobufany) ] |  | No |

#### v1beta1CreateNamespaceRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string |  | Yes |
| format | [SchemaFormat](#schemaformat) |  | No |
| compatibility | [SchemaCompatibility](#schemacompatibility) |  | No |
| description | string |  | No |

#### v1beta1CreateNamespaceResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| namespace | [v1beta1Namespace](#v1beta1namespace) |  | No |

#### v1beta1CreateSchemaResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| version | integer |  | No |
| id | string |  | No |
| location | string |  | No |

#### v1beta1DeleteNamespaceResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |

#### v1beta1DeleteSchemaResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |

#### v1beta1DeleteVersionResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |

#### v1beta1GetLatestSchemaResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| data | byte |  | No |

#### v1beta1GetNamespaceResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| namespace | [v1beta1Namespace](#v1beta1namespace) |  | No |

#### v1beta1GetSchemaMetadataResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| format | [SchemaFormat](#schemaformat) |  | No |
| compatibility | [SchemaCompatibility](#schemacompatibility) |  | No |
| authority | string |  | No |

#### v1beta1GetSchemaResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| data | byte |  | No |

#### v1beta1ListNamespacesResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| namespaces | [ string ] |  | No |

#### v1beta1ListSchemasResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| schemas | [ string ] |  | No |

#### v1beta1ListVersionsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| versions | [ integer ] |  | No |

#### v1beta1Namespace

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string |  | No |
| format | [SchemaFormat](#schemaformat) |  | No |
| Compatibility | [SchemaCompatibility](#schemacompatibility) |  | No |
| description | string |  | No |
| createdAt | dateTime |  | No |
| updatedAt | dateTime |  | No |

#### v1beta1SearchHits

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| namespaceId | string |  | No |
| schemaId | string |  | No |
| versionId | integer |  | No |
| fields | [ string ] |  | No |
| types | [ string ] |  | No |

#### v1beta1SearchMeta

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| total | long |  | No |

#### v1beta1SearchResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| hits | [ [v1beta1SearchHits](#v1beta1searchhits) ] |  | No |
| meta | [v1beta1SearchMeta](#v1beta1searchmeta) |  | No |

#### v1beta1UpdateNamespaceResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| namespace | [v1beta1Namespace](#v1beta1namespace) |  | No |

#### v1beta1UpdateSchemaMetadataResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| format | [SchemaFormat](#schemaformat) |  | No |
| compatibility | [SchemaCompatibility](#schemacompatibility) |  | No |
| authority | string |  | No |
