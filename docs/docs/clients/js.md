# JavaScript

Stencil nodejs client package provides a store to lookup protobuf descriptors and options to keep the protobuf descriptors upto date.

It has following features

- Ability to refresh protobuf descriptors in specified intervals
- Support to download descriptors from multiple urls

## Installation

```sh
npm install --save @raystack/stencil
```

## Usage

### Creating a client

```js
const { Stencil } = require("stencil");

const url = "http://url/to/proto/descriptorset/file";
const client = await Stencil.getInstance(url, {
  shouldRefresh: true,
  refreshInterval: 43200, // 12 hrs
});
```

### Creating a multiURLClient

```js
const { MultiURLStencil } = require("stencil");

const url1 = "http://url/to/proto/descriptorset/file";
const url2 = "http://url/to/proto/descriptorset/file2";
const client = await MultiURLStencil.getInstance([url1, url2], {
  shouldRefresh: true,
  refreshInterval: 43200, // 12 hrs
});
```

### Get proto descriptor type

```js
const { Stencil } = require("stencil");

const url = "http://url/to/proto/descriptorset/file";
const client = await Stencil.getInstance(url, {
  shouldRefresh: false,
});
const type = client.getType("google.protobuf.DescriptorProto");
```

### Encode/Decode message

Let's say we want to encode message for below proto message defniniton

```proto
syntax = "proto3";

package test;

message One {
  int64 field_one = 1;
}
```

```js
const { Stencil } = require('stencil');

const url = 'http://url/to/proto/descriptorset/file';
const client = await Stencil.getInstance(url, {
  shouldRefresh: false
});
const type = client.getType('test.One');
// Encode
const msg = { field_one: 10 };
const errs = type.verify(msg);
if errs {
   throw new Error(`unable to serialize message: ${errs}`);
}
const encodedBuffer = type.encode(msg).finish();
// Decode
const decodedType = type.decode(encodedBuffer);
console.log(decodedType.toObject())
```

## Setting up development environment

### Prerequisite Tools

- [Node.js](https://nodejs.org/) (version >= 12.0.0)
- [Git](https://git-scm.com/)

1. Clone the repo

   ```sh
   $ git clone https://github.com/raystack/stencil
   $ cd stencil/clients/js
   ```

2. Install dependencies

   ```sh
   $ npm install
   ```

3. Run the tests. All of the tests are written with [jest](https://jestjs.io/).

   ```sh
   $ npm test
   ```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags](https://github.com/raystack/stencil/tags).
