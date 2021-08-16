# Stencil Server

Stencil Server is written in Go. It provides REST interface for storing and retrieving protobuf descriptorset file.

## Features

 - stores versioned history of proto descriptor file on specified namespace and name
 - enforce backward compatability check on upload by default
 - ability to skip some of the backward compatability checks while upload
 - ability to download fully contained proto descriptor file for specified proto message [fullName](https://pkg.go.dev/google.golang.org/protobuf@v1.27.1/reflect/protoreflect#FullName)
 - provides metadata API to retrieve latest version number given a name and namespace

## Deployment

This section describes Deployment instructions for stencil server

{% page-ref page="deployment.md" %}

## Rules

This section explains all backward compatability rules available in stencil server

{% page-ref page="rules.md" %}

## API reference

This section contains complete API reference for stencil server

{% page-ref page="api.md" %}
