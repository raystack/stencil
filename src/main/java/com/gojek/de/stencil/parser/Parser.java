package com.gojek.de.stencil.parser;

import com.google.protobuf.DynamicMessage;
import com.google.protobuf.InvalidProtocolBufferException;

public interface Parser {
    DynamicMessage parse(byte[] data) throws InvalidProtocolBufferException;
}

