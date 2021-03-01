# Stencil

Stencil is dynamic schema registry for protobuf. Protobuf is a great efficient and fast mechanism for serializing structured data. The challenge with protobuf is that for every change it requires to recompile the package to generate the necessary classes. This is not a big challenge if you have protobuf enclosed in your application and compile at startup. But if you have thousands of protos stored in central registry and 100s of applications use them. Updating dependencies of compiled proto jar can soon become a nightmare.

Protobuf allows you to define a whole proto file using [google.protobuf.FileDescriptorProto](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L62). A [google.protobuf.FileDescriptorSet](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto#L57) contains list of FileDescriptorProto. Stencil heavily make use of this feature to update proto schemas in runtime.
