# Stencil Python client

[![PyPI version](https://badge.fury.io/py/stencil-python-client.svg)](https://pypi.org/project/stencil-python-client)

Stencil Python client package provides a store to lookup protobuf descriptors and options to keep the protobuf descriptors upto date.

It has following features
 - Deserialize protobuf messages directly by specifying protobuf message name
 - Ability to refresh protobuf descriptors in specified intervals
 - Support to download descriptors from multiple urls


## Requirements

 - Python 3.7+

## Installation

Use `pip`
```
pip3 install stencil-python-client
```

Then import the stencil package into your own code as mentioned below
```python
from raystack import stencil
```

## Usage

### Creating a client

```python
from raystack import stencil

url = "http://url/to/proto/descriptorset/file"
client = stencil.Client(url)
```

### Creating a multiURLClient

```python
from raystack import stencil

urls = ["http://urlA", "http://urlB"]
client = stencil.MultiUrlClient(urls)
```

### Get Descriptor
```python
from raystack import stencil

url = "http://url/to/proto/descriptorset/file"
client = stencil.Client(url)
client.get_descriptor("google.protobuf.DescriptorProto")
```

### Parse protobuf message. 
```python
from raystack import stencil

url = "http://url/to/proto/descriptorset/file"
client = stencil.Client(url)

data = ""
desc = client.parse("google.protobuf.DescriptorProto", data)
```

Refer to [stencil documentation](https://raystack.gitbook.io/stencil/) for more information what you can do in stencil.
