package io.odpf.stencil.cache;

import java.util.Map;

import com.google.protobuf.Descriptors;

public abstract class ProtoUpdateListener {
    private String proto;

    public ProtoUpdateListener(String proto) {
        this.proto = proto;
    }

    public String getProto() {
        return proto;
    }

    public abstract void onProtoUpdate(String url, final Map<String, Descriptors.Descriptor> newDescriptor);
}
