package io.odpf.stencil.client;

import com.google.protobuf.Descriptors;
import io.odpf.stencil.models.DescriptorAndTypeName;

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

    boolean shouldAutoRefreshCache();
}