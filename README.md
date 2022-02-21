# Stencil

![Test](https://github.com/odpf/stencil/actions/workflows/test-server.yaml/badge.svg)
![Release](https://github.com/odpf/stencil/actions/workflows/release-server.yml/badge.svg)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?logo=apache)](LICENSE)
[![Version](https://img.shields.io/github/v/release/odpf/stencil?logo=semantic-release)](Version)

Stencil is a schema registry that provides schema mangement and validation to ensure data compatibility across applications. It enables developers to create, manage and consume schemas dynamically, efficiently, and reliably, and provides a simple way to validate data against those schemas. Stencil support multiple formats including Protobuf, Avro and JSON.

<p align="center"><img src="./docs/static/assets/overview.svg" /></p>

## Key Features

Discover why users choose Stencil as their main schema registry

- **Version history** Stencil stores versioned history of proto descriptor file on specified namespace and name
- **Backward compatibility** enforce backward compatability check on upload by default
- **Flexbility** ability to skip some of the backward compatability checks while upload
- **Descriptor fetch** ability to download proto descriptor files
- **Metadata** provides metadata API to retrieve latest version number given a name and namespace
- **Clients in multiple languages** Stencil provides clients in GO, Java, JS languages to interact with Stencil server and deserialize messages using dynamic schema

## Documentation

Explore the following resources to get started with Stencil:

- [Documentation](http://odpf.github.io/stencil) provides guidance on using stencil.
- [Server](https://odpf.github.io/stencil/docs/server/overview) provides details on getting started with stencil server.
- [Clients](https://odpf.github.io/stencil/docs/clients/overview) provides reference to supported stencil clients.

## Installation

Install Stencil on macOS, Windows, Linux, OpenBSD, FreeBSD, and on any machine.

#### Binary (Cross-platform)

Download the appropriate version for your platform from [releases](https://github.com/odpf/stencil/releases) page. Once downloaded, the binary can be run from anywhere.
You don’t need to install it into a global location. This works well for shared hosts and other systems where you don’t have a privileged account.
Ideally, you should install it somewhere in your PATH for easy use. `/usr/local/bin` is the most probable location.

#### Homebrew

```sh
# Install stencil (requires homebrew installed)
$ brew install odpf/taps/stencil

# Upgrade stencil (requires homebrew installed)
$ brew upgrade stencil

# Check for installed stencil version
$ stencil version
```

## Usage

Stencil has three major components. Server, CLI and clients. Stencil server and CLI are bundled in a single binary.

**Server**

Stencil server provides a way to store and fetch schemas and enforce compatability rules. Run `stencil server --help` to see instructions to manage Stencil server.

**CLI**

Stencil CLI allows users to iteract with server to create, view, and search schemas. CLI is fully featured but simple to use, even for those who have very limited experience working from the command line. Run `stencil --help` to see list of all available commands and instructions to use.

**Clients**

Stencil clients allows application to interact with stencil server to eserialize and deserialize messages using schema. Stencil supports clients in multiple languages.

- [Java](clients/java)
- [Go](clients/go)
- [Javascript](clients/js)
- Ruby - Coming soon
- Python - Coming soon

## Running locally

<details>
  <summary>Dependencies:</summary>

    - Git
    - Go 1.16 or above
    - PostgreSQL 13 or above

</details>

```sh
# Clone the repo
$ git clone git@github.com:odpf/stencil.git

# Check all build comamnds available
$ make help

# Build meteor binary file
$ make build

# Init server config
$ cp app/config.yaml config.yaml

# Run database migrations
$ ./stencil server migrate

# Start stencil server
$ ./stencil server start
```

## Running tests

```sh
# Running all unit tests
$ make test

# Print code coverage
$ make coverage
```

## Contribute

Development of Stencil happens in the open on GitHub, and we are grateful to the community for contributing bugfixes and improvements. Read below to learn how you can take part in improving stencil.

Read our [contributing guide](docs/contribute/contribution.md) to learn about our development process, how to propose bugfixes and improvements, and how to build and test your changes to stencil.

To help you get your feet wet and get you familiar with our contribution process, we have a list of [good first issues](https://github.com/odpf/stencil/labels/good%20first%20issue) that contain bugs which have a relatively limited scope. This is a great place to get started.

This project exists thanks to all the [contributors](https://github.com/odpf/stencil/graphs/contributors).

## License

Stencil is [Apache 2.0](LICENSE) licensed.
