(defproject io.odpf/stencil-clj "0.2.0-SNAPSHOT"
  :description "Stencil client for clojure"
  :url "https://github.com/odpf/stencil"
  :license {:name "EPL-2.0 OR GPL-2.0-or-later WITH Classpath-exception-2.0"
            :url "https://www.eclipse.org/legal/epl-2.0/"}
  :dependencies [[org.clojure/clojure "1.10.3"]
                 [io.odpf/stencil "0.2.0"]]
  :plugins [[lein-cljfmt "0.7.0"]]
  :global-vars {*warn-on-reflection* true}
  :source-paths ["src"]
  :repl-options {:init-ns stencil.core})
