package io.odpf.stencil.parser;

import com.google.protobuf.AbstractMessage;
import com.google.protobuf.InvalidProtocolBufferException;

public interface Parser {
    AbstractMessage parse(byte[] data) throws InvalidProtocolBufferException;
}

