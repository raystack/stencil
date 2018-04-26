package com.gojek.de.stencil;

import com.google.protobuf.Descriptors;

public interface StencilClient {
    public Descriptors.Descriptor get(String className);
    public void load();
    public void reload();
}
