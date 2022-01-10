package io.odpf.stencil.client;

import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;
import com.google.protobuf.InvalidProtocolBufferException;
import io.odpf.stencil.exception.StencilRuntimeException;

import java.io.Closeable;
import java.util.Map;

/**
 * A client to get the protobuf descriptors and more information
 */
public interface StencilClient extends Closeable {
    Descriptors.Descriptor get(String className);

    default DynamicMessage parse(String className, byte[] data) throws InvalidProtocolBufferException {
        Descriptors.Descriptor descriptor = get(className);
        if (descriptor == null) {
            throw new StencilRuntimeException(new Throwable(String.format("No Descriptors found for %s", className)));
        }
        return DynamicMessage.parseFrom(descriptor, data);
    }

    Map<String, Descriptors.Descriptor> getAll();

    void refresh();
}
