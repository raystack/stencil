package io.odpf.stencil.cache;

import io.odpf.stencil.models.DescriptorAndTypeName;

import java.util.Map;

public abstract class ProtoUpdateListener {
    private String proto;

    public ProtoUpdateListener(String proto) {
        this.proto = proto;
    }

    public String getProto() {
        return proto;
    }

    public abstract void onProtoUpdate(String url, final Map<String, DescriptorAndTypeName> newDescriptor);
}
