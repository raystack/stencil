(ns stencil.decode
  (:require [clojure.string :as string])
  (:import [com.google.protobuf Descriptors$Descriptor Descriptors$FieldDescriptor Descriptors$EnumValueDescriptor DynamicMessage ByteString]))

(defn- byte-string->bytes
  [^ByteString value]
  (.toByteArray value))

(defn- replace-underscores-to-hyphen [k]
  (string/replace k #"_" "-"))

(defn- field-name->keyword
  [k]
  (-> k
      (replace-underscores-to-hyphen)
      (keyword)))

(defn- enum-value->clj-name
  [^Descriptors$EnumValueDescriptor value]
  (-> value
      (.getName)
      (field-name->keyword)))

(defn- fd->keyword
  [^Descriptors$FieldDescriptor fd]
  (let [name (.getName fd)]
    (field-name->keyword name)))

(declare proto-message->clojure-map)

(defn- proto-field->clojure-map-fn
  [^Descriptors$FieldDescriptor fd]
  (let [type-name (-> fd
                      (.getType)
                      (.toString))
        transform-fn (case type-name
                       ("INT32" "UINT32" "SINT32" "FIXED32" "SFIXED32") int
                       ("INT64" "UINT64" "SINT64" "FIXED64" "SFIXED64") long
                       "DOUBLE" double
                       "FLOAT" float
                       "BOOL" boolean
                       "STRING" str
                       "BYTES" byte-string->bytes
                       "ENUM" enum-value->clj-name
                       "MESSAGE" proto-message->clojure-map)]
    (if (.isRepeated fd)
      (partial map transform-fn)
      transform-fn)))

(defn- proto-message->clojure-map
  [^DynamicMessage msg]
  (let [all-fields (.getAllFields msg)
        reducer (fn [acc [k v]]
                  (->> ((proto-field->clojure-map-fn k) v)
                       (assoc acc (fd->keyword k))))]
    (reduce reducer {} all-fields)))

(defn- get-dynamic-message
  [^Descriptors$Descriptor desc ^"[B" data]
  (DynamicMessage/parseFrom desc data))

(defn bytes->map
  [^Descriptors$Descriptor desc ^"[B" data]
  (-> (get-dynamic-message desc data)
      proto-message->clojure-map))
