package com.gojek.de.stencil.client;

import com.google.protobuf.Descriptors;

import java.io.Closeable;

public interface StencilClient extends Closeable {
    Descriptors.Descriptor get(String className);

    default String getAppName() {
        String podName = System.getenv("POD_NAME");
        if (podName != null) return podName;
        return "";
    }
}