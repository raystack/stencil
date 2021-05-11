# Stencil Server

Stencil Server is written in Go. It provides REST interface for storing and retrieving protobuf descriptorset file.

## Features

 - stores versioned history of proto descriptor file on specified namespace and name
 - enforce backward compatability check on upload by default
 - ability to skip some of the backward compatability checks while upload
 - ability to download proto descriptor files
 - provides metadata API to retrieve latest version number given a name and namespace
 - ability to download latest proto descriptor file
 - support for multiple backend storage services (Local storage, Google cloud storage, S3, Azure blob storage and in-memory storage)

## Deployment

This section describes Deployment instructions for stencil server

{% page-ref page="deployment.md" %}

## Rules

This section explains all backward compatability rules available in stencil server

{% page-ref page="rules.md" %}

## API reference

This section contains complete API reference for stencil server

{% page-ref page="api.md" %}
