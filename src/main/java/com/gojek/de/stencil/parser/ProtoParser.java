package com.gojek.de.stencil.parser;

import com.gojek.de.stencil.client.StencilClient;
import com.gojek.de.stencil.exception.StencilRuntimeException;
import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;
import com.google.protobuf.InvalidProtocolBufferException;
import com.timgroup.statsd.StatsDClient;

import java.time.Instant;
import java.util.Optional;

public class ProtoParser implements Parser {
    private StencilClient stencilClient;
    private Optional<StatsDClient> statsDClientOpt;
    private String protoClassName;

    public ProtoParser(StencilClient stencilClient, Optional<StatsDClient> statsDClientOpt, String protoClassName) {
        this.stencilClient = stencilClient;
        this.statsDClientOpt = statsDClientOpt;
        this.protoClassName = protoClassName;
    }

    public DynamicMessage parse(byte[] bytes) throws InvalidProtocolBufferException {
        Instant start = Instant.now();
        Descriptors.Descriptor descriptor = stencilClient.get(protoClassName);
        Instant end = Instant.now();
        long latencyMillis = end.toEpochMilli() - start.toEpochMilli();
        statsDClientOpt.ifPresent(s -> s.recordExecutionTime("stencil.exec.time", latencyMillis, "name:name"));
        if (descriptor == null) {
            throw new StencilRuntimeException(new Throwable(String.format("No Descriptors found for %s", protoClassName)));
        }
        return DynamicMessage.parseFrom(descriptor, bytes);
    }
}
