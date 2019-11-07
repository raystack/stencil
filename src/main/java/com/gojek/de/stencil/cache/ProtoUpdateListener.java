package com.gojek.de.stencil.cache;

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
}
