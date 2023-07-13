(ns stencil.core
  (:require [stencil.encode :refer [map->bytes]]
            [stencil.decode :refer [bytes->map]])
  (:import
   (org.raystack.stencil.client StencilClient)
   (org.raystack.stencil StencilClientFactory)
   (org.raystack.stencil.cache SchemaRefreshStrategy)
   (org.raystack.stencil.exception StencilRuntimeException)
   (org.raystack.stencil.config StencilConfig)
   (org.apache.http.message BasicHeader) (java.util ArrayList)))

(defn create-client
  "Returns a new Stencil Clojure client instance by passing client-config.

       ### Client config structure :
       | Key                    | Type      | Description                                                                                 |
       | -----------------------|-----------|---------------------------------------------------------------------------------------------|
       | `url`                  | _String_  | Stencil url to fetch latest schema                                                          |
       | `refresh-cache`        | _Boolean_ | Whether the cache should be refreshed or not                                                |
       | `refresh-ttl`          | _Integer_ | Cache TTL in minutes                                                                        |
       | `request-timeout`      | _Integer_ | Request timeout in milliseconds                                                             |
       | `request-backoff-time` | _Integer_ | Request back off time in minutes                                                            |
       | `retry-count`          | _Integer_ | Number of retries to be made to fetch schema                                                |
       | `headers`              | _Map_     | Map with key as header key and value as header value, which will be passed to stencil server|
       | `refresh-strategy`     | _keyword_ | Possible values :version-based-refresh, :long-polling-refresh. Default :long-polling-refresh|

       Example:
       ```clojure
       (let [sample-client-config {:url       \"https://example-url\"
                                  :refresh-cache        true
                                  :refresh-ttl          100
                                  :request-timeout      10000
                                  :request-backoff-time 100
                                  :retry-count          3
                                  :headers              {\"Authorization\" \"Bearer <token>\"}
                                  :refresh-strategy     :version-based-refresh
                                  }]
             (create-client sample-client-config))
       ```"
  [client-config]
  (try (let [stencil-config (-> (StencilConfig/builder)
                                (.fetchTimeoutMs (int (:request-timeout client-config)))
                                (.fetchRetries (int (:retry-count client-config)))
                                (.cacheAutoRefresh (:refresh-cache client-config))
                                (.cacheTtlMs (long (:refresh-ttl client-config)))
                                (.fetchHeaders (let [array-list (new ArrayList)]
                                                 (doseq [[k v] (:headers client-config)]
                                                   (.add array-list (BasicHeader. k v)))
                                                 array-list))
                                (.refreshStrategy (if (= :version-based-refresh (:refresh-strategy client-config))
                                                    (SchemaRefreshStrategy/versionBasedRefresh)
                                                    (SchemaRefreshStrategy/longPollingStrategy)))
                                (.build))]
         (StencilClientFactory/getClient (:url client-config)
                                         stencil-config))
       (catch StencilRuntimeException e (throw (ex-info "Client initialization failed" {:cause :client-initialization-failed
                                                                                        :info  e})))))

(defn get-descriptor
  "Returns protobuf descriptor object for the given protobuf class name."
  [^StencilClient client proto-class-name]
  (.get client proto-class-name))

(defn deserialize
  "Returns Clojure map for the given protobuf encoded byte array and protobuf class name."
  [^StencilClient client proto-class-name data]
  (-> (get-descriptor client proto-class-name)
      (bytes->map data)))

(defn serialize
  "Returns protobuf encoded byte array for the given Clojure and protobuf class name."
  [^StencilClient client proto-class-name map]
  (-> (get-descriptor client proto-class-name)
      (map->bytes map)))

