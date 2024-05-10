# Clojure

A Clojure library designed to easily encode and decode protobuf messages by using Clojure maps.

## Installation

Add the below dependency to your `project.clj` file:

```clj
[com.gotocompany/stencil-clj "0.5.0"]
```

## Usage

Consider following proto message

```proto
syntax = "proto3";

package example;

option java_multiple_files = true;
option java_package = "com.goto.CljTest";

message Address {
	string city = 1;
	string street = 2;
}

message Person {
	enum Gender {
		UNKNOWN = 0;
		MALE = 1;
		FEMALE = 2;
		NON_BINARY = 3;
	}
	string name = 1;
	Address address = 2;
	Gender gender = 3;
	repeated string email_list = 4;
	int32 age = 5;
}
```

1. Create stencil client. You can refer to [java client](java) documentation for all available options.

```clojure
(ns test
  (:require [stencil.core :refer [create-client]]))

(def client (create-client {:url "<stencil service url>"
                :refresh-cache true
                :refresh-strategy :version-based-refresh
                :headers {"<headerkey>" "<header value>"}))
```

2. To serialize data from clojure map

```clojure
(:require [stencil.core :refer [serialize]])

(def serialized-data
     (serialize client "com.goto.CljTest" {:name "Foo"
                                          :address {:street "bar"}
                                          :email-list ["a@example.com" "b@b.com"]
                                          :gender :NON-BINARY
                                          :age 10}))
```

3. Deserialize data from bytes to clojure map

```clojure
(:require [stencil.core :refer [deserialize]])

(deserialize client "com.goto.CljTest" serialized-data)
;; prints
;; {:name "Foo"
;; :address {:street "bar"}
;; :email-list ["a@example.com" "b@b.com"]
;; :gender :NON-BINARY
;; :age 10}
```

## Protocol buffers - Clojure interop

| Protobuf        | Clojure                                                                                                                          | Notes                                                                                               |
| --------------- | -------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| field names     | keywords in kebab case                                                                                                           | `name` -> `:name`, `field_name` -> `:field-name`                                                    |
| scalar fields   | Values follow [protobuf-java scalar value mappings](https://developers.google.com/protocol-buffers/docs/proto3#scalar)           |                                                                                                     |
| enums           | Values converted as keywords of enum's original value                                                                            | `UNKNOWN` -> `:UNKNOWN`                                                                             |
| messages        | clojure map                                                                                                                      | `message Hello {string name = 1;}` -> {:name "goto"}                                                |
| repeated fields | clojure vector                                                                                                                   |                                                                                                     |
| one-of fields   | treated as regular fields                                                                                                        | if two fields are set that are part of one-of, last seen value is considered while serializing data |
| map             | map values follow it's [wire representation](https://developers.google.com/protocol-buffers/docs/proto3#backwards_compatibility) | for `map<string, string>` type, example value will be `[{:key "key" :value "value"}]`               |

**Note on errors:**
Serialize will throw error in following cases

1. unknown field is passed that's not present in schema `{:cause :unknown-field :info {:field-name <field-name>}}`
2. if non-collection type is passed to repeated field `{:cause :not-a-collection :info {:value <value>}}`
3. If unknown enum value passed that's not present in schema `{:cause :unknown-enum-value :info {:field-name <field-name>}}`

## API

- `create-client (client-config)`

  Returns a new Stencil Clojure client instance by passing client-config.

  ### Client config structure :

  | Key                    | Type      | Description                                                                                  |
  | ---------------------- | --------- | -------------------------------------------------------------------------------------------- |
  | `url`                  | _String_  | Stencil url to fetch latest descriptor sets                                                  |
  | `refresh-cache`        | _Boolean_ | Whether the cache should be refreshed or not                                                 |
  | `refresh-ttl`          | _Integer_ | Cache TTL in minutes                                                                         |
  | `request-timeout`      | _Integer_ | Request timeout in milliseconds                                                              |
  | `request-backoff-time` | _Integer_ | Request back off time in minutes                                                             |
  | `retry-count`          | _Integer_ | Number of retries to be made to fetch descriptor sets                                        |
  | `headers`              | _Map_     | Map with key as header key and value as header value, which will be passed to stencil server |
  | `refresh-strategy`     | _keyword_ | Possible values :version-based-refresh, :long-polling-refresh. Default :long-polling-refresh |

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

  | Key                | Type     | Description                                                       |
  | ------------------ | -------- | ----------------------------------------------------------------- |
  | `client`           | _Object_ | Instantiated Clojure client object                                |
  | `proto-class-name` | _String_ | Name of the proto class whose proto descriptor object is required |

  ### Response structure

  | Value          | Type     | Description                                    |
  | -------------- | -------- | ---------------------------------------------- |
  | **proto-desc** | _Object_ | Protobuf descriptor for given proto class name |

  Example:

  ```clojure
  (let [client (create-client sample-client-config)
        proto-package "com.goto.stencil_clj_test"
        proto-class-name "Scalar"
        fully-qualified-proto-name (str proto-package "." proto-class-name)]
      (get-descriptor client fully-qualified-proto-name))
  ```

- `deserialize (client proto-class-name data)`

  Returns Clojure map for the given protobuf encoded byte array and protobuf class name.

  ### Argument list :

  | Key                | Type         | Description                                                        |
  | ------------------ | ------------ | ------------------------------------------------------------------ |
  | `client`           | _Object_     | Instantiated Clojure client object                                 |
  | `proto-class-name` | _String_     | Name of the proto class whose proto descriptor object is required  |
  | `data`             | _Byte-Array_ | Data (byte-array) to be deserialized using proto-descriptor object |

  ### Response structure

  | Value                    | Type                 | Description                        |
  | ------------------------ | -------------------- | ---------------------------------- |
  | **deserialized-message** | _PersistentArrayMap_ | Deserialized message (Clojure Map) |

  Example:

  ```clojure
  (let [client (create-client sample-client-config)
        proto-package "com.goto.stencil_clj_test"
        proto-class-name "Scalar"
        fully-qualified-proto-name (str proto-package "." proto-class-name)
        proto-desc (get-descriptor client fully-qualified-proto-name)
        data-to-deserialize (serialize client fully-qualified-proto-name{:field-one 1.25})]
       (deserialize client fully-qualified-proto-name data-to-deserialize))
  ```

- `serialize (client proto-class-name map)`

  Returns protobuf encoded byte array for the given Clojure and protobuf class name.

  ### Argument list :

  | Key                | Type                 | Description                                                              |
  | ------------------ | -------------------- | ------------------------------------------------------------------------ |
  | `client`           | _Object_             | Instantiated Clojure client object                                       |
  | `proto-class-name` | _String_             | Name of the proto class whose proto descriptor object is required        |
  | `map`              | _PersistentArrayMap_ | Data (in the form of map) to be serialized using proto descriptor object |

  ### Response structure

  | Value                  | Type         | Description                     |
  | ---------------------- | ------------ | ------------------------------- |
  | **serialized-message** | _Byte-Array_ | Serialized message (byte-array) |

  Example:

  ```clojure
  (let [client (create-client sample-client-config)
        proto-package "com.goto.stencil_clj_test"
        proto-class-name "Scalar"
        fully-qualified-proto-name (str proto-package "." proto-class-name)
        proto-desc (get-descriptor client fully-qualified-proto-name)]
       (serialize client fully-qualified-proto-name {:field-one 1.25}))
  ```

## Development

- Ensure [leiningen](https://leiningen.org/) is installed.

- Run tests: `lein clean && lein javac && lein test`

- Run formatting: `lein cljfmt fix`
