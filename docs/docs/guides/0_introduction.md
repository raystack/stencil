import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";

# Introduction

This tour introduces you to Stencil schema registry. Along the way you will learn how to manage schemas, enforce rules, serialise and deserialise data using stencil clients.

### Prerequisites

This tour requires you to have Stencil CLI tool installed on your local machine. You can run `stencil version` to verify the installation. Please follow installation guide if you do not have it installed already.

Stencil CLI and clients talks to Stencil server to publish and fetch schema. Please make sure you also have a stencil server running. You can also run server locally with `stencil server start` command. For more details check deployment guide.

### Help

At any time you can run the following commands.

```bash
# Check the installed version for stencil cli tool
$ stencil version

# See the help for a command
$ stencil --help
```

Help command can also be run on any sub command.

```bash
$ stencil schema --help
```

Check the reference for stencil cli commands.

```bash
$ stencil reference
```
