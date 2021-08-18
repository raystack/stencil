<h1 id="stencil-server">Stencil server v0.1.4</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

<h1 id="stencil-server-default">Default</h1>

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

<h1 id="stencil-server-descriptors">descriptors</h1>

Manage descriptors

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
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Success response|None|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict|None|

<aside class="success">
This operation does not require authentication
</aside>

## get__v1_namespaces_{namespace}_descriptors

> Code samples

```shell
# You can also use wget
curl -X GET /v1/namespaces/{namespace}/descriptors \
  -H 'Accept: application/json'

```

`GET /v1/namespaces/{namespace}/descriptors`

*list all available descriptor names under one namespace*

<h3 id="get__v1_namespaces_{namespace}_descriptors-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|namespace|path|string|true|none|

> Example responses

> 200 Response

```json
[
  "string"
]
```

<h3 id="get__v1_namespaces_{namespace}_descriptors-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|returns list of descriptor names|Inline|

<h3 id="get__v1_namespaces_{namespace}_descriptors-responseschema">Response Schema</h3>

<aside class="success">
This operation does not require authentication
</aside>

## get__v1_namespaces_{namespace}_descriptors_{name}_versions

> Code samples

```shell
# You can also use wget
curl -X GET /v1/namespaces/{namespace}/descriptors/{name}/versions \
  -H 'Accept: application/json'

```

`GET /v1/namespaces/{namespace}/descriptors/{name}/versions`

*list all available versions for specified descriptor*

<h3 id="get__v1_namespaces_{namespace}_descriptors_{name}_versions-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|namespace|path|string|true|none|
|name|path|string|true|none|

> Example responses

> 200 Response

```json
[
  "string"
]
```

<h3 id="get__v1_namespaces_{namespace}_descriptors_{name}_versions-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|returns list of versions|Inline|

<h3 id="get__v1_namespaces_{namespace}_descriptors_{name}_versions-responseschema">Response Schema</h3>

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

<h1 id="stencil-server-metadata">metadata</h1>

manage latest versions for uploaded descriptor files

## post__v1_namespaces_{namespace}_metadata

> Code samples

```shell
# You can also use wget
curl -X POST /v1/namespaces/{namespace}/metadata \
  -H 'Content-Type: application/json'

```

`POST /v1/namespaces/{namespace}/metadata`

*update metadata*

> Body parameter

```json
{
  "name": "string",
  "version": "string"
}
```

<h3 id="post__v1_namespaces_{namespace}_metadata-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|namespace|path|string|true|none|
|body|body|[MetadataPayload](#schemametadatapayload)|true|specify name and version in payload|

<h3 id="post__v1_namespaces_{namespace}_metadata-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Success response|None|

<aside class="success">
This operation does not require authentication
</aside>

## get__v1_namespaces_{namespace}_metadata_{name}

> Code samples

```shell
# You can also use wget
curl -X GET /v1/namespaces/{namespace}/metadata/{name} \
  -H 'Accept: application/json'

```

`GET /v1/namespaces/{namespace}/metadata/{name}`

*get latest version for specified descriptor*

<h3 id="get__v1_namespaces_{namespace}_metadata_{name}-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|namespace|path|string|true|none|
|name|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "version": "string"
}
```

<h3 id="get__v1_namespaces_{namespace}_metadata_{name}-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Success response|[MetadataResponse](#schemametadataresponse)|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_MetadataResponse">MetadataResponse</h2>
<!-- backwards compatibility -->
<a id="schemametadataresponse"></a>
<a id="schema_MetadataResponse"></a>
<a id="tocSmetadataresponse"></a>
<a id="tocsmetadataresponse"></a>

```json
{
  "version": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|version|string|false|none|none|

<h2 id="tocS_MetadataPayload">MetadataPayload</h2>
<!-- backwards compatibility -->
<a id="schemametadatapayload"></a>
<a id="schema_MetadataPayload"></a>
<a id="tocSmetadatapayload"></a>
<a id="tocsmetadatapayload"></a>

```json
{
  "name": "string",
  "version": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|false|none|none|
|version|string|false|none|none|




