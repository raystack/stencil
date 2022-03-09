(ns stencil.core-test
  (:require [clojure.test :refer :all]
            [stencil.core :refer :all])
  (:import
   (io.odpf.stencil.client StencilClient)))

(deftest test-create-client
  (testing "should create client"
    (let [config {:url "http://localhost:8000/v1beta1/namespaces/odpf/schemas/proton"
                  :refresh-ttl          100
                  :request-timeout      10000
                  :request-backoff-time 100
                  :retry-count          3
                  :refresh-cache true
                  :headers {"Authorization" "Bearer token"}
                  :refresh-strategy :long-polling-refresh}
          client (create-client config)]
      (is (instance? StencilClient client)))))
