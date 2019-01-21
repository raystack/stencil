package com.gojek.de.stencil.parser;

import com.gojek.de.stencil.client.StencilClient;
import com.gojek.de.stencil.exception.StencilRuntimeException;
import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;
import com.google.protobuf.InvalidProtocolBufferException;
import com.timgroup.statsd.NoOpStatsDClient;
import com.timgroup.statsd.StatsDClient;

import java.time.Instant;

public class ProtoParser implements Parser {
    private StencilClient stencilClient;
    private StatsDClient statsDClient;
    private String protoClassName;

    public ProtoParser(StencilClient stencilClient, StatsDClient statsDClient, String protoClassName) {
        this.stencilClient = stencilClient;
        this.statsDClient = statsDClient;
        this.protoClassName = protoClassName;
    }

    public ProtoParser(StencilClient stencilClient, String protoClassName) {
        this(stencilClient, new NoOpStatsDClient(), protoClassName);
    }

    public DynamicMessage parse(byte[] bytes) throws InvalidProtocolBufferException {
        Instant start = Instant.now();
        Descriptors.Descriptor descriptor = stencilClient.get(protoClassName);
        Instant end = Instant.now();
        long latencyMillis = end.toEpochMilli() - start.toEpochMilli();
        statsDClient.recordExecutionTime("stencil.exec.time", latencyMillis, "name=" + stencilClient.getAppName());
        if (descriptor == null) {
            throw new StencilRuntimeException(new Throwable(String.format("No Descriptors found for %s", protoClassName)));
        }
        return DynamicMessage.parseFrom(descriptor, bytes);
    }
}
