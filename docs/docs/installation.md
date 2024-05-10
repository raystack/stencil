# Installation

Stencil installation is simple. You can install Stencil on macOS, Windows, Linux, OpenBSD, FreeBSD, and on any machine. There are several approaches to installing Stencil.

1. Using a [pre-compiled binary](#binary-cross-platform)
3. Installing with [Docker](#Docker)
4. Installing from [source](#building-from-source)

### Binary (Cross-platform)

Download the appropriate version for your platform from [releases](https://github.com/goto/stencil/releases) page. Once downloaded, the binary can be run from anywhere.
You don’t need to install it into a global location. This works well for shared hosts and other systems where you don’t have a privileged account.
Ideally, you should install it somewhere in your `PATH` for easy use. `/usr/local/bin` is the most probable location.

### Docker

We provide ready to use Docker container images. To pull the latest image:

```
docker pull gotocompany/stencil:latest
```

To pull a specific version:

```
docker pull gotocompany/stencil:0.8.1
```

### Building from source

To compile from source, you will need [Go](https://golang.org/) installed and a copy of [git](https://www.git-scm.com/) in your `PATH`.

```bash
# Clone the repo
$ git clone git@github.com:goto/stencil.git

# Check all build comamnds available
$ make help

# Build stencil binary file
$ make build

# Check for installed stencil version
$ ./stencil version
```

### Verifying the installation

To verify Stencil is properly installed, run `stencil --help` on your system. You should see help output. If you are executing it from the command line, make sure it is on your `PATH` or you may get an error about Stencil not being found.

```bash
$ stencil --help
```
