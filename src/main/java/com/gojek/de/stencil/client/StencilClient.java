package com.gojek.de.stencil.client;

import com.gojek.de.stencil.models.DescriptorAndTypeName;
import com.google.protobuf.Descriptors;

import java.io.Closeable;
import java.util.Map;

public interface StencilClient extends Closeable {
    Descriptors.Descriptor get(String className);

    Map<String, Descriptors.Descriptor> getAll();

    Map<String, String> getTypeNameToPackageNameMap();

    Map<String, DescriptorAndTypeName> getAllDescriptorAndTypeName();

    default String getAppName() {
        String podName = System.getenv("POD_NAME");
        if (podName != null) return podName;
        return "";
    }

    void refresh();
}