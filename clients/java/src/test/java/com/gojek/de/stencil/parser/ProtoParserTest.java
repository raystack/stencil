package com.gojek.de.stencil.parser;

import com.gojek.de.stencil.DescriptorMapBuilder;
import com.gojek.de.stencil.client.StencilClient;
import com.gojek.de.stencil.exception.StencilRuntimeException;
import com.gojek.de.stencil.models.DescriptorAndTypeName;
import com.gojek.stencil.TestMessage;
import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;
import com.google.protobuf.InvalidProtocolBufferException;
import org.junit.Before;
import org.junit.Test;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Map;

import static org.junit.Assert.assertNotNull;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;


public class ProtoParserTest {

    private static final String LOOKUP_KEY = "com.gojek.stencil.TestMessage";
    Map<String, DescriptorAndTypeName> descriptorMap;

    @Before
    public void setup() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String DESCRIPTORS_FILE_PATH = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTORS_FILE_PATH).getFile());
        descriptorMap = new DescriptorMapBuilder().buildFrom(fileInputStream);
    }

    @Test(expected = StencilRuntimeException.class)
    public void testProtoParseOnException() throws InvalidProtocolBufferException {
        StencilClient stencilClient = mock(StencilClient.class);
        ProtoParser protoParser = new ProtoParser(stencilClient, LOOKUP_KEY);
        protoParser.parse(TestMessage.getDefaultInstance().toByteArray());
    }

    @Test
    public void testProtoParseOnSuccess() throws InvalidProtocolBufferException {
        StencilClient stencilClient = mock(StencilClient.class);
        when(stencilClient.get(LOOKUP_KEY)).thenReturn(descriptorMap.get(LOOKUP_KEY).getDescriptor());
        ProtoParser protoParser = new ProtoParser(stencilClient, LOOKUP_KEY);
        DynamicMessage parsed = protoParser.parse(TestMessage.newBuilder().setSampleString("sample_string").build().toByteArray());
        assertNotNull(parsed);
    }
}
