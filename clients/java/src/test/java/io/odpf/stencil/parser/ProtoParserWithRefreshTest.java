package io.odpf.stencil.parser;

import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;
import com.google.protobuf.InvalidProtocolBufferException;
import io.odpf.stencil.DescriptorMapBuilder;
import io.odpf.stencil.TestMessage;
import io.odpf.stencil.TestMessageSuperset;
import io.odpf.stencil.client.StencilClient;
import io.odpf.stencil.exception.StencilRuntimeException;
import io.odpf.stencil.models.DescriptorAndTypeName;
import org.apache.commons.codec.DecoderException;
import org.junit.Before;
import org.junit.Test;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Map;

import static org.junit.Assert.assertNotNull;
import static org.mockito.Mockito.*;


public class ProtoParserWithRefreshTest {

    private static final String LOOKUP_KEY = "io.odpf.stencil.TestMessage";
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

        byte[] testData = TestMessageSuperset.newBuilder().setSampleString("sample_string").setSuccess(true).build().toByteArray();
        DynamicMessage parsed = protoParser.parse(testData);
        assertNotNull(parsed);
        verify(stencilClient, times(1)).refresh();
    }
}
