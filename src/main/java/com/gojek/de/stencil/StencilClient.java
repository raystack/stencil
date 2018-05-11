package com.gojek.de.stencil;

import com.google.protobuf.Descriptors;

public abstract class StencilClient {
    private static StencilClient instance;
    public abstract Descriptors.Descriptor get(String className);
    public abstract void load();
    public abstract void reload();
}
