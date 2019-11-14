package com.gojek.de.stencil.client;

import com.google.protobuf.Descriptors;

import java.io.Closeable;
import java.util.Map;

public interface StencilClient extends Closeable {
    Descriptors.Descriptor get(String className);
    Map<String, Descriptors.Descriptor> getAll();

    default String getAppName() {
        String podName = System.getenv("POD_NAME");
        if (podName != null) return podName;
        return "";
    }
}