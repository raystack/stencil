package com.gojek.de.stencil.cache;

import com.gojek.de.stencil.models.DescriptorAndTypeName;

import java.util.Map;

public abstract class ProtoUpdateListener {
    private String proto;

    public ProtoUpdateListener(String proto) {
        System.out.println(proto);
        this.proto = proto;
    }

    public String getProto() {
        return proto;
    }

    public abstract void onProtoUpdate();

    public abstract void onProtoUpdate(String url, final Map<String, DescriptorAndTypeName> newDescriptor);
}
