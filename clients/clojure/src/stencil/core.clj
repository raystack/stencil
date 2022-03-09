(ns stencil.core
  (:require [stencil.encode :refer [map->bytes]]
            [stencil.decode :refer [bytes->map]])
  (:import
   (io.odpf.stencil.client StencilClient)
   (io.odpf.stencil StencilClientFactory)
   (io.odpf.stencil.cache SchemaRefreshStrategy)
   (io.odpf.stencil.exception StencilRuntimeException)
   (io.odpf.stencil.config StencilConfig)
   (org.apache.http.message BasicHeader)))

(defn create-client
  "convenience function to create stencil java client instance. Users can use stencil java API directly to create stencil client instance.

   ### Client config structure :
   | Key                    | Type      | Description                                                                                |
   | -----------------------|-----------|--------------------------------------------------------------------------------------------|
   | `url`                  | _String_  | Stencil url                                                                                |
   | `refresh-cache`        | _Boolean_ | Whether the cache should be refreshed or not                                               |
   | `refresh-ttl`          | _Integer_ | Cache TTL in minutes                                                                       |
   | `request-timeout`      | _Integer_ | Request timeout in milliseconds                                                            |
   | `request-backoff-time` | _Integer_ | Request back off time in minutes                                                           |
   | `retry-count`          | _Integer_ | Number of retries to be made to fetch descriptor sets                                      |
   | `headers`              | _Map_     | Map with key as header key and value as header value, will be passed to stencil server     |
   | `refresh-strategy`     | _keyword_ | Possible values :version-based-refresh, :long-polling-refresh. Default:long-polling-refresh|

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
                                (.fetchHeaders (let [array-list (new java.util.ArrayList)]
                                                 (doseq [[k v] (:headers client-config)]
                                                   (.add array-list (BasicHeader. k v)))))
                                (.refreshStrategy (if (= :version-based-refresh (:refresh-strategy client-config))
                                                    (SchemaRefreshStrategy/versionBasedRefresh)
                                                    (SchemaRefreshStrategy/longPollingStrategy)))
                                (.build))]
         (StencilClientFactory/getClient (:url client-config)
                                         stencil-config))
       (catch StencilRuntimeException e (throw (ex-info "Client initialization failed" {:cause :client-initialization-failed
                                                                                        :info  e})))))

(defn get-descriptor
  "returns protobuf descriptor given a stencil client and classname"
  [^StencilClient client class-name]
  (.get client class-name))

(defn deserialize
  "returns clojure map for given protobuf encoded byte data and class name"
  [^StencilClient client class-name data]
  (-> (get-descriptor client class-name)
      (bytes->map data)))

(defn serialize
  "returns byte array for given clojure map and class name"
  [^StencilClient client class-name map]
  (-> (get-descriptor client class-name)
      (map->bytes map)))

