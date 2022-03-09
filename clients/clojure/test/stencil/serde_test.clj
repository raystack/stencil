(ns stencil.serde-test
  (:require [clojure.test :refer :all]
            [stencil.encode :refer [map->bytes]]
            [stencil.decode :refer [bytes->map]])
  (:import
   (io.odpf.stencil DescriptorMapBuilder)
   (java.io File FileInputStream)))

(defn file-desc-map [^String path]
  (let [f (File. path)
        is (FileInputStream. f)
        descriptor-map (DescriptorMapBuilder/buildFrom is)]
    (fn [key] (.get descriptor-map key))))

(def local-get-descriptor (file-desc-map "./test/stencil/testdata/testdata.desc"))

(def scalar-data {:field-one 1.25
                  :float-field 1.5
                  :field-int32 1234
                  :field-int64 1234567890123456789
                  :field-uint32 1234
                  :field-uint64 2345678901234567890
                  :field-sint32 9012
                  :field-sint64 3456789012345678901
                  :field-fixed32 4567
                  :field-fixed64 4567890123456789012
                  :field-sfixed32 8912
                  :field-sfixed64 5678901234567890123
                  :field-bool true
                  :field-string "abc"
                  :field-bytes (.getBytes "foo")})

(def simple-nested {:field-name :VALUE3
                    :group "abc"
                    :nested-field {:field-bool true
                                   :field-int32 1234}
                    :duration-field  {:seconds 1631789198
                                      :nanos   763000000}
                    :timestamp-field {:seconds 1631789198
                                      :nanos   763000000}})

(def complex-types {:name         "general"
                    :map-field    [{:key "key_1" :value {:field-one 1.24}}
                                   {:key "key-2" :value {:field-bool true}}]
                    :struct-field {:fields
                                   [{:key "num", :value {:number-value 1.0}}
                                    {:key "structkey", :value {:struct-value {:fields [{:key "num", :value {:number-value 3.0}}]}}}
                                    {:key "arraykey", :value {:list-value {:values [{:number-value 2.0}]}}}]}})

(def simple-array {:groups        [:VALUE-1 :VALUE-2 :VALUE-1]
                   :values        ["a" "b" "c"]
                   :nested-fields [{:field-name :VALUE-2
                                    :group      "group-val"}
                                   {:group "abc"}
                                   simple-nested]})
(def recursive-data {:name         "level-1"
                     :single-field {:name         "level-2"
                                    :single-field {:name        "level-3"
                                                   :multi-field [{:name "level-4-0"}
                                                                 {:name        "level-4-1"
                                                                  :group-field :VALUE-1}
                                                                 {:name         "level-4-2"
                                                                  :single-field {:name        "level-5"
                                                                                 :group-field :value4}}]}}})

(def wrapper-data {:one {:value "abc"}
                   :two {:value 2.0}
                   :three {:value 3.0}
                   :four {:value 4}
                   :five {:value 5}
                   :six {:value 6}
                   :seven {:value 7}
                   :eight {:value true}})

(defn- verify
  [name data]
  (let [descriptor (local-get-descriptor name)
        serialized-data (map->bytes descriptor data)
        deserialized-data (bytes->map descriptor serialized-data)]
    (is (= deserialized-data data))))

(deftest serialization-deserialization-test
  (testing "should handle scalar types"
    (let [proto-name "io.odpf.stencil_clj_test.Scalar"
          descriptor (local-get-descriptor proto-name)
          serialized-data (map->bytes descriptor scalar-data)
          deserialized-data (bytes->map descriptor serialized-data)]
      (is (= (dissoc deserialized-data :field-bytes) (dissoc scalar-data :field-bytes)))
      (is (= (seq (:field-bytes deserialized-data)) (seq (:field-bytes scalar-data))))))

  (testing "should handle enum type if enum value is by name"
    (verify "io.odpf.stencil_clj_test.SimpleNested" {:field-name :VALUE-1}))

  (testing "should handle enum type if enum value is by number"
    (let [proto-name "io.odpf.stencil_clj_test.SimpleNested"
          descriptor (local-get-descriptor proto-name)
          test-data {:field-name 2}
          serialized-data (map->bytes descriptor test-data)
          deserialized-data (bytes->map descriptor serialized-data)]
      (is (= deserialized-data {:field-name :VALUE-2}))))

  (testing "should deserialize message type field"
    (verify "io.odpf.stencil_clj_test.SimpleNested" simple-nested))

  (testing "should handle struct and map types"
    (verify "io.odpf.stencil_clj_test.ComplexTypes" complex-types))

  (testing "should handle repeated fields"
    (verify "io.odpf.stencil_clj_test.SimpleArray" simple-array))

  (testing "should handle self referencing types"
    (verify "io.odpf.stencil_clj_test.Recursive" recursive-data))

  (testing "should handle wrapper types"
    (verify "io.odpf.stencil_clj_test.Wrappers" wrapper-data)))
