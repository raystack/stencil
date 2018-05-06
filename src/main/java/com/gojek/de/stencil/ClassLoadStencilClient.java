package com.gojek.de.stencil;

import com.google.protobuf.Descriptors;

import java.util.HashMap;
import java.util.Map;

public class ClassLoadStencilClient extends StencilClient {

    private Map<String, Descriptors.Descriptor> descriptorMap;

    public ClassLoadStencilClient() {
        descriptorMap = new HashMap<>();
    }
    @Override
    public Descriptors.Descriptor get(String className) {
        if (! descriptorMap.containsKey(className)) {
            try {
                Class<?> protoClass = Class.forName(className);
                descriptorMap.put(className, (Descriptors.Descriptor) protoClass.getMethod("getDescriptor").invoke(null));
            } catch (ReflectiveOperationException exception) {

            }
        }
        return descriptorMap.get(className);
    }

    @Override
    public void load() {
        descriptorMap = new HashMap<>();
    }

    @Override
    public void reload() {
        load();
    }
}
