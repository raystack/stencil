# Installation

Stencil installation is simple. You can install Stencil on macOS, Windows, Linux, OpenBSD, FreeBSD, and on any machine. There are several approaches to installing Stencil.

1. Using a [pre-compiled binary](#binary-cross-platform)
2. Installing with [package manager](#MacOS)
3. Installing with [Docker](#Docker)
4. Installing from [source](#building-from-source)

### Binary (Cross-platform)

Download the appropriate version for your platform from [releases](https://github.com/raystack/stencil/releases) page. Once downloaded, the binary can be run from anywhere.
You don’t need to install it into a global location. This works well for shared hosts and other systems where you don’t have a privileged account.
Ideally, you should install it somewhere in your `PATH` for easy use. `/usr/local/bin` is the most probable location.

### MacOS

`stencil` is available via a Homebrew Tap, and as downloadable binary from the [releases](https://github.com/raystack/stencil/releases/latest) page:

```sh
brew install raystack/tap/stencil
```

To upgrade to the latest version:

```
brew upgrade stencil
```

#### Linux

`stencil` is available as downloadable binaries from the [releases](https://github.com/raystack/stencil/releases/latest) page. Download the `.deb` or `.rpm` from the releases page and install with `sudo dpkg -i` and `sudo rpm -i` respectively.

### Windows

`stencil` is available via [scoop](https://scoop.sh/), and as a downloadable binary from the [releases](https://github.com/raystack/stencil/releases/latest) page:

```
scoop bucket add stencil https://github.com/raystack/scoop-bucket.git
```

To upgrade to the latest version:

```
scoop update stencil
```

### Docker

We provide ready to use Docker container images. To pull the latest image:

```
docker pull raystack/stencil:latest
```

To pull a specific version:

```
docker pull raystack/stencil:v0.5.0
```

### Building from source

To compile from source, you will need [Go](https://golang.org/) installed and a copy of [git](https://www.git-scm.com/) in your `PATH`.

```bash
# Clone the repo
$ git clone git@github.com:raystack/stencil.git

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
