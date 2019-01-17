package com.gojek.de.stencil.client;

import com.google.protobuf.Descriptors;

import java.util.concurrent.ExecutionException;

public interface StencilClient {
    Descriptors.Descriptor get(String className);
}
