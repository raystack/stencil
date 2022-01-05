package io.odpf.stencil.client;

import com.google.common.testing.FakeTicker;
import com.google.protobuf.Descriptors;
import com.google.protobuf.DynamicMessage;
import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.Struct;
import io.odpf.stencil.DescriptorMapBuilder;
import io.odpf.stencil.NestedField;
import io.odpf.stencil.account_db_accounts.FULLDOCUMENT;
import io.odpf.stencil.cache.SchemaCacheLoader;
import io.odpf.stencil.config.StencilConfig;
import io.odpf.stencil.exception.StencilRuntimeException;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Map;
import java.util.concurrent.TimeUnit;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.mockito.Mockito.*;

public class URLStencilClientWithCacheTest {

    private Map<String, Descriptors.Descriptor> descriptorMap;
    private StencilClient stencilClient;
    private static final String DESCRIPTOR_FILE_PATH = "__files/descriptors.bin";
    private static final String LOOKUP_KEY = "io.odpf.stencil.TestMessage";

    @Before
    public void setup() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        descriptorMap = DescriptorMapBuilder.buildFrom(fileInputStream);
    }

    @After
    public void close() throws IOException {
        stencilClient.close();
    }

    @Test
    public void getFromStencilClientSuccessfully() {
        SchemaCacheLoader cacheLoader = mock(SchemaCacheLoader.class);
        when(cacheLoader.load(LOOKUP_KEY)).thenReturn(descriptorMap);

        stencilClient = new URLStencilClient(LOOKUP_KEY, StencilConfig.builder().build(), cacheLoader);
        Descriptors.Descriptor result = stencilClient.get(LOOKUP_KEY);

        verify(cacheLoader, times(1)).load(LOOKUP_KEY);
        verify(cacheLoader, times(0)).reload(LOOKUP_KEY, descriptorMap);
        assertNotNull(result);
    }

    @Test(expected = StencilRuntimeException.class)
    public void getFromStencilClientOnException() throws Exception {
        SchemaCacheLoader cacheLoader = mock(SchemaCacheLoader.class);
        when(cacheLoader.load(LOOKUP_KEY)).thenThrow(new StencilRuntimeException(new Throwable()));

        stencilClient = new URLStencilClient(LOOKUP_KEY, StencilConfig.builder().build(), cacheLoader);
        stencilClient.get(LOOKUP_KEY);
    }

    @Test
    public void shouldNotRefreshCacheByDefault() {
        SchemaCacheLoader cacheLoader = mock(SchemaCacheLoader.class);
        when(cacheLoader.load(LOOKUP_KEY)).thenReturn(descriptorMap);

        FakeTicker fakeTicker = new FakeTicker();

        stencilClient = new URLStencilClient(LOOKUP_KEY, StencilConfig.builder().build(), cacheLoader, fakeTicker);
        Descriptors.Descriptor result = stencilClient.get(LOOKUP_KEY);

        fakeTicker.advance(1000, TimeUnit.MINUTES);
        Descriptors.Descriptor reloadedResult = stencilClient.get(LOOKUP_KEY);

        verify(cacheLoader, times(1)).load(LOOKUP_KEY);
        verify(cacheLoader, times(0)).reload(LOOKUP_KEY, descriptorMap);
        assertEquals(result, reloadedResult);
    }

    @Test
    public void getFromStencilClientSuccessfullySubsequentTimesFromCache() {
        SchemaCacheLoader cacheLoader = mock(SchemaCacheLoader.class);
        when(cacheLoader.load(LOOKUP_KEY)).thenReturn(descriptorMap);

        FakeTicker fakeTicker = new FakeTicker();
        StencilConfig stencilConfig = StencilConfig.builder().cacheAutoRefresh(true).build();

        stencilClient = new URLStencilClient(LOOKUP_KEY, stencilConfig, cacheLoader, fakeTicker);
        Descriptors.Descriptor result = stencilClient.get(LOOKUP_KEY);
        assertNotNull(result);

        fakeTicker.advance(stencilConfig.getCacheTtlMs() + 1000, TimeUnit.MILLISECONDS);
        Descriptors.Descriptor reloadedResult = stencilClient.get(LOOKUP_KEY);

        verify(cacheLoader, times(1)).load(LOOKUP_KEY);
        verify(cacheLoader, times(1)).reload(LOOKUP_KEY, descriptorMap);
        assertNotNull(reloadedResult);
    }

    @Test(expected = StencilRuntimeException.class)
    public void shouldThrowExceptionOnDescriptorNotFound() throws InvalidProtocolBufferException {
        SchemaCacheLoader cacheLoader = mock(SchemaCacheLoader.class);
        when(cacheLoader.load(LOOKUP_KEY)).thenReturn(descriptorMap);

        FakeTicker fakeTicker = new FakeTicker();
        StencilConfig stencilConfig = StencilConfig.builder().cacheAutoRefresh(true).build();

        stencilClient = new URLStencilClient(LOOKUP_KEY, stencilConfig, cacheLoader, fakeTicker);
        NestedField msg = NestedField.newBuilder().setStringField("stencil").setIntField(10).build();
        stencilClient.parse("io.odpf.stencil.invalid", msg.toByteArray());
    }

    @Test(expected = InvalidProtocolBufferException.class)
    public void shouldThrowExceptionOnParsingInvalidData() throws InvalidProtocolBufferException {
        SchemaCacheLoader cacheLoader = mock(SchemaCacheLoader.class);
        descriptorMap.put("com.google.protobuf.Struct", Struct.getDescriptor());
        when(cacheLoader.load(LOOKUP_KEY)).thenReturn(descriptorMap);

        FakeTicker fakeTicker = new FakeTicker();
        StencilConfig stencilConfig = StencilConfig.builder().cacheAutoRefresh(true).build();

        stencilClient = new URLStencilClient(LOOKUP_KEY, stencilConfig, cacheLoader, fakeTicker);
        FULLDOCUMENT msg = FULLDOCUMENT.newBuilder().setCif("cifvalue").setId("idvalue").build();
        stencilClient.parse("com.google.protobuf.Struct", msg.toByteArray());
    }

    @Test
    public void parseShouldCreateDynamicMessageOnValidData() throws InvalidProtocolBufferException {
        SchemaCacheLoader cacheLoader = mock(SchemaCacheLoader.class);
        when(cacheLoader.load(LOOKUP_KEY)).thenReturn(descriptorMap);

        FakeTicker fakeTicker = new FakeTicker();
        StencilConfig stencilConfig = StencilConfig.builder().cacheAutoRefresh(true).build();

        stencilClient = new URLStencilClient(LOOKUP_KEY, stencilConfig, cacheLoader, fakeTicker);
        FULLDOCUMENT msg = FULLDOCUMENT.newBuilder().setCif("cifvalue").setId("idvalue").build();
        DynamicMessage newMsg = stencilClient.parse("io.odpf.stencil.account_db_accounts.FULLDOCUMENT", msg.toByteArray());
        Descriptors.Descriptor desc = stencilClient.get("io.odpf.stencil.account_db_accounts.FULLDOCUMENT");
        assertNotNull(newMsg);
        Object value = newMsg.getField(desc.findFieldByNumber(1));
        assertEquals("idvalue", value);
    }

}
