package com.gojek.de.stencil.parser;

import com.gojek.de.stencil.client.StencilClient;
import com.gojek.de.stencil.exception.StencilRuntimeException;
import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;
import com.google.protobuf.InvalidProtocolBufferException;
import com.timgroup.statsd.NoOpStatsDClient;
import com.timgroup.statsd.StatsDClient;

import java.time.Instant;

public class ProtoParserWithRefresh implements Parser {
    private StencilClient stencilClient;
    private StatsDClient statsDClient;
    private String protoClassName;

    public ProtoParserWithRefresh(StencilClient stencilClient, StatsDClient statsDClient, String protoClassName) {
        this.stencilClient = stencilClient;
        this.statsDClient = statsDClient;
        this.protoClassName = protoClassName;
    }

    public ProtoParserWithRefresh(StencilClient stencilClient, String protoClassName) {
        this(stencilClient, new NoOpStatsDClient(), protoClassName);
    }

    public DynamicMessage parse(byte[] bytes) throws InvalidProtocolBufferException {
        Instant start = Instant.now();
        Descriptors.Descriptor descriptor = getDescriptor();
        if (descriptor == null) {
            throw new StencilRuntimeException(new Throwable(String.format("No Descriptors found for %s", protoClassName)));
        }
        DynamicMessage parsedMessage = getMessage(bytes, descriptor);
        statsDClient.recordExecutionTime("stencil.exec.time,name=" + stencilClient.getAppName(), Instant.now().toEpochMilli() - start.toEpochMilli() );
        return parsedMessage;
    }

    private DynamicMessage getMessage(byte[] bytes, Descriptors.Descriptor descriptor) throws InvalidProtocolBufferException {
        DynamicMessage parsedMessage = DynamicMessage.parseFrom(descriptor, bytes);
        if (!hasUnknownFields(parsedMessage)) {
            return parsedMessage;
        }
        stencilClient.refresh();
        return DynamicMessage.parseFrom(descriptor, bytes);
    }

    private Descriptors.Descriptor getDescriptor() {
        Descriptors.Descriptor descriptor = stencilClient.get(protoClassName);
        if (descriptor != null) {
            return descriptor;
        }
        stencilClient.refresh();
        return stencilClient.get(protoClassName);
    }

    private boolean hasUnknownFields(DynamicMessage parsedMessage) {
        return parsedMessage.getUnknownFields().asMap().size() > 0;
    }
}
