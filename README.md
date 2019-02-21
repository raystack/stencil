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

### Publish stencil to central maven
In order to publish to central maven, you require sonatype credentials and [GnuPG](http://gnupg.org) setup.

To get sonatype credentials, Register to `https://issues.sonatype.org/`  
To setup GnuPG:
    Install gpg with `brew install gpg`
    Refer to this [url](https://docs.gradle.org/current/userguide/signing_plugin.html#sec:signatory_credentials) to setup GnuPG    

After this, add these values to `gradle.properties` in your user directory
```

signing.keyId=<last eight symbols of gnupg keyId> 
signing.password=<your passphrase to unlock gpg secrets>
signing.secretKeyRingFile=/Users/me/.gnupg/secring.gpg

ossrhUsername=your-jira-id
ossrhPassword=your-jira-password
```

Upload your gpg keys to ubuntu opengpg server (Required once)
Run the following command: 
```
gpg --keyserver hkp://pool.sks-keyservers.net --recv-keys <last eight symbols of gnupg keyId>
gpg --keyserver hkp://keyserver.ubuntu.com --send-keys <last eight symbols of gnupg keyId>
```

### Notes
- we used ``java-statsd-client`` from ``com.timgroup``, so use the same client in the application that integrates with stencil client
