<h1 id="stencil-server">Stencil server v0.1.4</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

<h1 id="stencil-server-health">health</h1>

## ping

<a id="opIdping"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /ping

```

`GET /ping`

*service health check*

<h3 id="ping-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|returns pong message|None|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="stencil-server-stencilservice">StencilService</h1>

## post__v1_namespaces_{namespace}_descriptors

> Code samples

```shell
# You can also use wget
curl -X POST /v1/namespaces/{namespace}/descriptors \
  -H 'Content-Type: multipart/form-data'

```

`POST /v1/namespaces/{namespace}/descriptors`

*upload descriptors*

> Body parameter

```yaml
name: string
version: string
latest: true
dryrun: true
skiprules:
  - FILE_NO_BREAKING_CHANGE
file: string

```

<h3 id="post__v1_namespaces_{namespace}_descriptors-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|namespace|path|string|true|none|
|body|body|object|true|none|
|» name|body|string|true|none|
|» version|body|string|true|version number for descriptor file. This should follow semantic version compatability|
|» latest|body|boolean|false|mark this descriptor file as latest|
|» dryrun|body|boolean|false|flag for dryRun|
|» skiprules|body|[string]|false|list of rules to skip|
|» file|body|string(binary)|true|descriptorset file to upload|

#### Enumerated Values

|Parameter|Value|
|---|---|
|» skiprules|FILE_NO_BREAKING_CHANGE|
|» skiprules|MESSAGE_NO_DELETE|
|» skiprules|FIELD_NO_BREAKING_CHANGE|
|» skiprules|ENUM_NO_BREAKING_CHANGE|

<h3 id="post__v1_namespaces_{namespace}_descriptors-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Success response if operation succeded|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation error response when user payload has missing required fields or currently being uploaded file is not backward compatible with previously uploaded file|None|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|conflict error reponse if namespace, name and version combination already present|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected internal error reponse|None|

<aside class="success">
This operation does not require authentication
</aside>

## get__v1_namespaces_{namespace}_descriptors_{name}_versions_{version}

> Code samples

```shell
# You can also use wget
curl -X GET /v1/namespaces/{namespace}/descriptors/{name}/versions/{version}

```

`GET /v1/namespaces/{namespace}/descriptors/{name}/versions/{version}`

*download specified descriptor file*

<h3 id="get__v1_namespaces_{namespace}_descriptors_{name}_versions_{version}-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|namespace|path|string|true|none|
|name|path|string|true|none|
|version|path|string|true|none|
|fullnames|query|array[string]|false|Proto fullnames|

<h3 id="get__v1_namespaces_{namespace}_descriptors_{name}_versions_{version}-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|download response|None|

<aside class="success">
This operation does not require authentication
</aside>

## StencilService_ListSnapshots

<a id="opIdStencilService_ListSnapshots"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/snapshots \
  -H 'Accept: application/json'

```

`GET /v1/snapshots`

<h3 id="stencilservice_listsnapshots-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|namespace|query|string|false|none|
|name|query|string|false|none|
|version|query|string|false|none|
|latest|query|boolean|false|none|

> Example responses

> 200 Response

```json
{
  "snapshots": [
    {
      "id": "string",
      "namespace": "string",
      "name": "string",
      "version": "string",
      "latest": true
    }
  ]
}
```

<h3 id="stencilservice_listsnapshots-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListSnapshotsResponse](#schemav1listsnapshotsresponse)|
|default|Default|An unexpected error response.|[rpcStatus](#schemarpcstatus)|

<aside class="success">
This operation does not require authentication
</aside>

## StencilService_PromoteSnapshot

<a id="opIdStencilService_PromoteSnapshot"></a>

> Code samples

```shell
# You can also use wget
curl -X PATCH /v1/snapshots/{id}/promote \
  -H 'Accept: application/json'

```

`PATCH /v1/snapshots/{id}/promote`

*PromoteSnapshot promotes particular snapshot version as latest*

<h3 id="stencilservice_promotesnapshot-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string(int64)|true|none|

> Example responses

> 200 Response

```json
{
  "snapshot": {
    "id": "string",
    "namespace": "string",
    "name": "string",
    "version": "string",
    "latest": true
  }
}
```

<h3 id="stencilservice_promotesnapshot-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1PromoteSnapshotResponse](#schemav1promotesnapshotresponse)|
|default|Default|An unexpected error response.|[rpcStatus](#schemarpcstatus)|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_protobufAny">protobufAny</h2>
<!-- backwards compatibility -->
<a id="schemaprotobufany"></a>
<a id="schema_protobufAny"></a>
<a id="tocSprotobufany"></a>
<a id="tocsprotobufany"></a>

```json
{
  "typeUrl": "string",
  "value": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|typeUrl|string|false|none|none|
|value|string(byte)|false|none|none|

<h2 id="tocS_rpcStatus">rpcStatus</h2>
<!-- backwards compatibility -->
<a id="schemarpcstatus"></a>
<a id="schema_rpcStatus"></a>
<a id="tocSrpcstatus"></a>
<a id="tocsrpcstatus"></a>

```json
{
  "code": 0,
  "message": "string",
  "details": [
    {
      "typeUrl": "string",
      "value": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|code|integer(int32)|false|none|none|
|message|string|false|none|none|
|details|[[protobufAny](#schemaprotobufany)]|false|none|none|

<h2 id="tocS_v1ListSnapshotsResponse">v1ListSnapshotsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listsnapshotsresponse"></a>
<a id="schema_v1ListSnapshotsResponse"></a>
<a id="tocSv1listsnapshotsresponse"></a>
<a id="tocsv1listsnapshotsresponse"></a>

```json
{
  "snapshots": [
    {
      "id": "string",
      "namespace": "string",
      "name": "string",
      "version": "string",
      "latest": true
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|snapshots|[[v1Snapshot](#schemav1snapshot)]|false|none|none|

<h2 id="tocS_v1PromoteSnapshotResponse">v1PromoteSnapshotResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1promotesnapshotresponse"></a>
<a id="schema_v1PromoteSnapshotResponse"></a>
<a id="tocSv1promotesnapshotresponse"></a>
<a id="tocsv1promotesnapshotresponse"></a>

```json
{
  "snapshot": {
    "id": "string",
    "namespace": "string",
    "name": "string",
    "version": "string",
    "latest": true
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|snapshot|[v1Snapshot](#schemav1snapshot)|false|none|none|

<h2 id="tocS_v1Snapshot">v1Snapshot</h2>
<!-- backwards compatibility -->
<a id="schemav1snapshot"></a>
<a id="schema_v1Snapshot"></a>
<a id="tocSv1snapshot"></a>
<a id="tocsv1snapshot"></a>

```json
{
  "id": "string",
  "namespace": "string",
  "name": "string",
  "version": "string",
  "latest": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string(int64)|false|none|none|
|namespace|string|true|none|none|
|name|string|false|none|none|
|version|string|false|none|none|
|latest|boolean|false|none|none|

