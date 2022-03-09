(ns stencil.encode
  (:import [com.google.protobuf Descriptors$Descriptor Descriptors$FieldDescriptor Descriptors$EnumDescriptor DynamicMessage DynamicMessage$Builder ByteString])
  (:require [clojure.string :as string]))

(defn- replace-hyphen-to-underscores [k]
  (string/replace k #"-" "_"))

(defn- bytes->byte-string
  [^"[B" value]
  (ByteString/copyFrom value))

(defn- keyword->field-name
  [k]
  (-> k
      replace-hyphen-to-underscores
      (string/replace #"^:" "")
      name))

(defn- clj-name->enum-value
  [^Descriptors$FieldDescriptor fd value]
  (let [^Descriptors$EnumDescriptor ed (.getEnumType fd)]
    (if (number? value)
      (.findValueByNumber ed value)
      (.findValueByName ed (keyword->field-name value)))))

(declare clojure-map->proto-message)

(defn- clojure-map->proto-field-fn
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
                       "BYTES" bytes->byte-string
                       "ENUM" (partial clj-name->enum-value fd)
                       "MESSAGE" (partial clojure-map->proto-message (.getMessageType fd)))]
    (if (.isRepeated fd)
      (partial mapv transform-fn)
      transform-fn)))

(defn- keyword->fd
  [^Descriptors$Descriptor desc key]
  (->> key
       keyword->field-name
       (.findFieldByName desc)))

(defn- proto-message-builder
  [^Descriptors$Descriptor desc]
  (fn [^DynamicMessage$Builder acc [k v]]
    (let [^Descriptors$FieldDescriptor fd (keyword->fd desc k)]
      (if (nil? fd)
        acc
        (as-> fd $
          (clojure-map->proto-field-fn $)
          ($ v)
          (.setField acc fd $))))))

(defn- clojure-map->proto-message
  [^Descriptors$Descriptor desc data]
  (let [builder (DynamicMessage/newBuilder desc)]
    (-> desc
        proto-message-builder
        (reduce builder data)
        .build)))

(defn map->bytes
  [^Descriptors$Descriptor desc data]
  (-> (clojure-map->proto-message desc data)
      .toByteArray))

