package com.gojek.de.stencil.models;

import com.google.protobuf.Descriptors;

public class DescriptorAndTypeName {
    private Descriptors.Descriptor descriptor;
    private String typeName;

    public DescriptorAndTypeName(Descriptors.Descriptor descriptor, String typeName) {
        this.descriptor = descriptor;
        this.typeName = typeName;
    }

    public Descriptors.Descriptor getDescriptor() {
        return descriptor;
    }

    public String getTypeName() {
        return typeName;
    }
}
