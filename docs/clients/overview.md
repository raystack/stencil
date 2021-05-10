# Stencil Clients

Stencil clients abstracts handling of descriptorset file on client side. Currently we officially support Stencil client in Java, Go, JS languages.

## Features

- downloading of descriptorset file from server
- parse API to deserialize protobuf encoded messages
- lookup API to find proto descriptors
- inbuilt strategies to refresh protobuf schema definitions.

## A note on configuring Stencil clients

- Stencil server provides API to download latest descriptor file. If new version is available latest file will point to new descriptor file. Always use latest version proto descriptor url for stencil client if you want to refresh schema definitions in runtime.
- Keep the refresh intervals relatively large (eg: 24hrs or 12 hrs) to reduce the number of calls. It's unlikely that proto changes reflected across systems within a hour or so.
- You can refresh descriptor file only if unknowns fields are faced by the client while parsing. This reduces unneccessary frequent calls made by clients. Currently this feature supported in JAVA and GO clients.

## Java client

This section describes Java client documentation

{% page-ref page="java.md" %}

## Go Client

This section describes GO client documentation

{% page-ref page="go.md" %}

## JS client

This section describes JS client documentation

{% page-ref page="js.md" %}
