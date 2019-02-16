# Stencil client

To get the proto descriptors from a server

### Usage

### Gradle
`compile group: 'com.gojek.de', name: 'stencil-client', version: '2.0.0' `


### Properties (default value)

```
STENCIL_TIMEOUT_MS (10000)
STENCIL_BACKOFF_MS (2000-4000)
STENCIL_RETRIES (4)
TTL_IN_MINUTES (30-60)
```


### Example (pseudo code)
```

 StencilClient stencilClient = StencilClientFactory.getClient(appConfig.getStencilUrl(),
 new HashMap<>(), // all the user defined properties
 getStatsDClient());
 ProtoParser protoParser = new ProtoParser(stencilClient, appConfig.getProtoSchema());
 protoParser.parse(byte[])

```

### Notes
- we used ``java-statsd-client`` from ``com.timgroup``, so use the same client in the application that integrates with stencil client
