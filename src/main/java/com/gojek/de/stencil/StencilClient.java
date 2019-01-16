package com.gojek.de.stencil;

import com.google.protobuf.Descriptors;

import java.util.concurrent.ExecutionException;

public abstract class StencilClient {
    public abstract Descriptors.Descriptor get(String className);
}
