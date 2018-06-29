package com.gojek.de.stencil;

import com.google.protobuf.Descriptors;

public abstract class StencilClient {
    public abstract Descriptors.Descriptor get(String className);
}
