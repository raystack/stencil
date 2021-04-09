package io.odpf.stencil.parser;

import io.odpf.stencil.client.StencilClient;
import io.odpf.stencil.exception.StencilRuntimeException;
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
    private Descriptors.Descriptor descriptor;

    public ProtoParserWithRefresh(StencilClient stencilClient, StatsDClient statsDClient, String protoClassName) {
        if (stencilClient.shouldAutoRefreshCache()) {
            throw new UnsupportedOperationException(String.format("REFRESH_CACHE is not supported with %s", getClass().getName()));
        }
        this.stencilClient = stencilClient;
        this.statsDClient = statsDClient;
        this.protoClassName = protoClassName;
        this.descriptor = getDescriptor();
    }

    public ProtoParserWithRefresh(StencilClient stencilClient, String protoClassName) {
        this(stencilClient, new NoOpStatsDClient(), protoClassName);
    }

    public DynamicMessage parse(byte[] bytes) throws InvalidProtocolBufferException {
        Instant start = Instant.now();
        DynamicMessage parsedMessage = DynamicMessage.parseFrom(descriptor, bytes);
        if (hasUnknownFields(parsedMessage)) {
            parsedMessage = DynamicMessage.parseFrom(getRefreshedDescriptor(), bytes);
        }
        statsDClient.recordExecutionTime("stencil.exec.time,name=" + stencilClient.getAppName(), Instant.now().toEpochMilli() - start.toEpochMilli() );
        return parsedMessage;
    }

    private Descriptors.Descriptor getRefreshedDescriptor() {
        stencilClient.refresh();
        return getDescriptor();
    }

    private Descriptors.Descriptor getDescriptor() {
        descriptor = stencilClient.get(protoClassName);
        if (descriptor == null) {
            throw new StencilRuntimeException(new Throwable(String.format("No Descriptors found for %s", protoClassName)));
        }
        return descriptor;
    }

    private boolean hasUnknownFields(DynamicMessage parsedMessage) {
        return parsedMessage.getUnknownFields().asMap().size() > 0;
    }
}
