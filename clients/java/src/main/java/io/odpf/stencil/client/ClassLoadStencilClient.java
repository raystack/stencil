package io.odpf.stencil.client;

import com.google.protobuf.Descriptors;

import java.io.Serializable;
import java.util.HashMap;
import java.util.Map;

/**
 * {@link StencilClient} implementation that can fetch descriptor from Protobuf Descritor classes in classpath
 */
public class ClassLoadStencilClient implements Serializable, StencilClient {

    transient private Map<String, Descriptors.Descriptor> descriptorMap;

    public ClassLoadStencilClient() {
    }

    @Override
    public Descriptors.Descriptor get(String className) {
        if (descriptorMap == null) {
            descriptorMap = new HashMap<>();
        }
        if (!descriptorMap.containsKey(className)) {
            try {
                Class<?> protoClass = Class.forName(className);
                descriptorMap.put(className, (Descriptors.Descriptor) protoClass.getMethod("getDescriptor").invoke(null));
            } catch (ReflectiveOperationException ignored) {

            }
        }
        return descriptorMap.get(className);
    }

    @Override
    public Map<String, Descriptors.Descriptor> getAll() {
        throw new UnsupportedOperationException();
    }

    @Override
    public void close() {
    }

    @Override
    public void refresh() {
        throw new UnsupportedOperationException();
    }
}
