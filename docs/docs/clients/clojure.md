# Clojure

A Clojure library designed to easily encode and decode protobuf messages by using Clojure maps.

## Installation

Add the below dependency to your `project.clj` file:
```clj
           [io.odpf/stencil-clj "0.2.0-SNAPSHOT"]
```

## Usage
- `create-client (client-config)`

  Returns a new Stencil Clojure client instance by passing client-config.

  ### Client config structure :
  | Key                    | Type      | Description                                                                                 |
  | -----------------------|-----------|---------------------------------------------------------------------------------------------|
  | `url`                  | _String_  | Stencil url to fetch latest descriptor sets                                                 |
  | `refresh-cache`        | _Boolean_ | Whether the cache should be refreshed or not                                                |
  | `refresh-ttl`          | _Integer_ | Cache TTL in minutes                                                                        |
  | `request-timeout`      | _Integer_ | Request timeout in milliseconds                                                             |
  | `request-backoff-time` | _Integer_ | Request back off time in minutes                                                            |
  | `retry-count`          | _Integer_ | Number of retries to be made to fetch descriptor sets                                       |
  | `headers`              | _Map_     | Map with key as header key and value as header value, which will be passed to stencil server|
  | `refresh-strategy`     | _keyword_ | Possible values :version-based-refresh, :long-polling-refresh. Default :long-polling-refresh|

  Example:
  ```clojure
   (let [sample-client-config {:url       "https://example-url"
                              :refresh-cache        true
                              :refresh-ttl          100
                              :request-timeout      10000
                              :request-backoff-time 100
                              :retry-count          3
                              :headers              {"Authorization" "Bearer <token>"}
                              :refresh-strategy     :version-based-refresh
                              }]
         (create-client sample-client-config))
  ```

- `get-descriptor (client proto-class-name)`

  Returns protobuf descriptor object for the given protobuf class name.

  ### Argument list :
  | Key                                             | Type              | Description                                                                 |
  | ------------------------------------------------|-------------------|-----------------------------------------------------------------------------|
  | `client`                                        | _Object_          | Instantiated Clojure client object                                          |
  | `proto-class-name`                              | _String_          | Name of the proto class whose proto descriptor object is required           |

  ### Response structure
  | Value                                           | Type              | Description                                                                 |
  |-------------------------------------------------|-------------------|-----------------------------------------------------------------------------|
  | **proto-desc**                                  | _Object_          | Protobuf descriptor for given proto class name                              |

  Example:
  ```clojure
  (let [client (create-client sample-client-config)
        proto-package "io.odpf.stencil_clj_test"
        proto-class-name "Scalar"
        fully-qualified-proto-name (str proto-package "." proto-class-name)]
      (get-descriptor client fully-qualified-proto-name))
  ```

- `deserialize (client proto-class-name data)`

  Returns Clojure map for the given protobuf encoded byte array and protobuf class name.

  ### Argument list :
  | Key                                             | Type                | Description                                                                 |
  | ------------------------------------------------|---------------------|-----------------------------------------------------------------------------|
  | `client`                                        | _Object_            | Instantiated Clojure client object                                          |
  | `proto-class-name`                              | _String_            | Name of the proto class whose proto descriptor object is required           |
  | `data`                                          | _Byte-Array_        | Data (byte-array) to be deserialized using proto-descriptor object          |

  ### Response structure
  | Value                                           | Type                | Description                                                                 |
  |-------------------------------------------------|---------------------|-----------------------------------------------------------------------------|
  | **deserialized-message**                        | _PersistentArrayMap_| Deserialized message (Clojure Map)                                          |

  Example:
  ```clojure
  (let [client (create-client sample-client-config)
        proto-package "io.odpf.stencil_clj_test"
        proto-class-name "Scalar"
        fully-qualified-proto-name (str proto-package "." proto-class-name)
        proto-desc (get-descriptor client fully-qualified-proto-name)
        data-to-deserialize (serialize client fully-qualified-proto-name{:field-one 1.25})]
       (deserialize client fully-qualified-proto-name data-to-deserialize))
  ```

- `serialize (client proto-class-name map)`

  Returns protobuf encoded byte array for the given Clojure and protobuf class name.

  ### Argument list :
  | Key                                             | Type                 | Description                                                                 |
  | ------------------------------------------------|----------------------|-----------------------------------------------------------------------------|
  | `client`                                        | _Object_             | Instantiated Clojure client object                                          |
  | `proto-class-name`                              | _String_             | Name of the proto class whose proto descriptor object is required           |
  | `map`                                           | _PersistentArrayMap_ | Data (in the form of map) to be serialized using proto descriptor object    |

  ### Response structure
  | Value                                           | Type                | Description                                                                  |
  |-------------------------------------------------|---------------------|------------------------------------------------------------------------------|
  | **serialized-message**                          | _Byte-Array_        | Serialized message (byte-array)                                              |

  Example:
  ```clojure
  (let [client (create-client sample-client-config)
        proto-package "io.odpf.stencil_clj_test"
        proto-class-name "Scalar"
        fully-qualified-proto-name (str proto-package "." proto-class-name)
        proto-desc (get-descriptor client fully-qualified-proto-name)]
       (serialize client fully-qualified-proto-name {:field-one 1.25}))
  ```
## Development
- Ensure [leiningen](https://leiningen.org/) is installed.

- Run tests: ```lein clean && lein javac && lein test```

- Run formatting: ```lein cljfmt fix```


## License

Stencil clojure client is [Apache 2.0](LICENSE) licensed.