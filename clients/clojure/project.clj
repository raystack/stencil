(defproject io.odpf/stencil-clj "0.2.0-SNAPSHOT"
  :description "Stencil client for clojure"
  :url "https://github.com/odpf/stencil"
  :license {:name "Apache 2.0"
            :url "https://www.apache.org/licenses/LICENSE-2.0"}
  :dependencies [[org.clojure/clojure "1.10.3"]
                 [io.odpf/stencil "0.2.0"]]
  :global-vars {*warn-on-reflection* true}
  :source-paths ["src"]
  :repl-options {:init-ns stencil.core})
