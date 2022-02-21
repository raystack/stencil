# Installation

Stencil installation is simple. You can install Stencil on macOS, Windows, Linux, OpenBSD, FreeBSD, and on any machine. There are several approaches to installing Stencil.

1. Using a [pre-compiled binary](#binary-cross-platform)
2. Installing with [package manager](#homebrew)
3. Installing from [source](#building-from-source)
4. Installing with [Docker](#using-docker-image)

### Binary (Cross-platform)

Download the appropriate version for your platform from [releases](https://github.com/odpf/stencil/releases) page. Once downloaded, the binary can be run from anywhere.
You don’t need to install it into a global location. This works well for shared hosts and other systems where you don’t have a privileged account.
Ideally, you should install it somewhere in your `PATH` for easy use. `/usr/local/bin` is the most probable location.

### Homebrew

You can install `stencil` on macOS or Linux using Homebrew:

```bash
# Install stencil (requires homebrew installed)
$ brew install odpf/taps/stencil

# Upgrade stencil (requires homebrew installed)
$ brew upgrade stencil

# Check for installed stencil version
$ stencil version
```

### Building from source

To compile from source, you will need [Go](https://golang.org/) installed and a copy of [git](https://www.git-scm.com/) in your `PATH`.

```bash
# Clone the repo
$ git clone git@github.com:odpf/stencil.git

# Check all build comamnds available
$ make help

# Build stencil binary file
$ make build

# Check for installed stencil version
$ ./stencil version
```

### Using Docker image

Stencil ships a Docker image [odpf/stencil](https://hub.docker.com/r/odpf/stencil) that enables you to use `stencil` as part of your Docker workflow.

For example, you can run `stencil help` with this command:

```bash
$ docker run odpf/stencil --help
```

### Verifying the installation

To verify Stencil is properly installed, run `stencil --help` on your system. You should see help output. If you are executing it from the command line, make sure it is on your `PATH` or you may get an error about Stencil not being found.

```bash
$ stencil --help
```
