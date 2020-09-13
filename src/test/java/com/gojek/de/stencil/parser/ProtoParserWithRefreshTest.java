package com.gojek.de.stencil.parser;

import com.gojek.de.stencil.DescriptorMapBuilder;
import com.gojek.de.stencil.client.StencilClient;
import com.gojek.de.stencil.exception.StencilRuntimeException;
import com.gojek.de.stencil.models.DescriptorAndTypeName;
import com.gojek.stencil.TestMessage;
import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;
import com.google.protobuf.InvalidProtocolBufferException;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.junit.Before;
import org.junit.Test;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Map;

import static org.junit.Assert.assertNotNull;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.times;


public class ProtoParserWithRefreshTest {

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
        ProtoParserWithRefresh protoParser = new ProtoParserWithRefresh(stencilClient, LOOKUP_KEY);
        protoParser.parse(TestMessage.getDefaultInstance().toByteArray());
        verify(stencilClient, times(1)).refresh();
    }

    @Test
    public void testProtoParseWithReload() throws IOException, DecoderException {
        StencilClient stencilClient = mock(StencilClient.class);
        when(stencilClient.get(LOOKUP_KEY)).thenReturn(descriptorMap.get(LOOKUP_KEY).getDescriptor());
        ProtoParserWithRefresh protoParser = new ProtoParserWithRefresh(stencilClient, LOOKUP_KEY);

        //simulate TestMessage.newBuilder().setSampleString("sample_string").setSampleSecondString("second").build().toByteArray();
        byte[] testData = Hex.decodeHex("0a0d73616d706c655f737472696e6712067365636f6e64".toCharArray());

        DynamicMessage parsed = protoParser.parse(testData);
        assertNotNull(parsed);
        verify(stencilClient, times(1)).refresh();
    }
}
