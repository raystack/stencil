(defproject com.gotocompany/stencil-clj "0.5.0"
  :description "Stencil client for clojure"
  :url "https://github.com/goto/stencil"
  :license {:name "Apache 2.0"
            :url "https://www.apache.org/licenses/LICENSE-2.0"}
  :dependencies [[org.clojure/clojure "1.10.3"]
                 [com.gotocompany/stencil "0.7.0"]]
  :plugins [[lein-cljfmt "0.7.0"]]
  :global-vars {*warn-on-reflection* true}
  :source-paths ["src"]
  :repl-options {:init-ns stencil.core})
