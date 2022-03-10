# Overview

Stencil clients abstracts handling of descriptorset file on client side. Currently we officially support Stencil client in Java, Go, JS languages.

## Features

- downloading of descriptorset file from server
- parse API to deserialize protobuf encoded messages
- lookup API to find proto descriptors
- inbuilt strategies to refresh protobuf schema definitions.

## A note on configuring Stencil clients

- Stencil server provides API to download latest descriptor file. If new version is available latest file will point to new descriptor file. Always use latest version proto descriptor url for stencil client if you want to refresh schema definitions in runtime.
- Keep the refresh intervals relatively large (eg: 24hrs or 12 hrs) to reduce the number of calls depending on how fast systems produce new messages using new proto schema.
- You can refresh descriptor file only if unknowns fields are faced by the client while parsing. This reduces unneccessary frequent calls made by clients. Currently this feature supported in JAVA and GO clients.

## Languages

- [Java](java)
- [Go](go)
- [Javascript](js)
- [Clojure](clojure)
- Ruby - Coming soon
- Python - Coming soon
