package io.odpf.stencil.client;

import io.odpf.stencil.DescriptorMapBuilder;
import io.odpf.stencil.cache.DescriptorCacheLoader;
import io.odpf.stencil.config.StencilConfig;
import io.odpf.stencil.exception.StencilRuntimeException;
import io.odpf.stencil.models.DescriptorAndTypeName;
import com.google.common.testing.FakeTicker;
import com.google.protobuf.Descriptors;
import org.aeonbits.owner.ConfigFactory;
import org.junit.Before;
import org.junit.Test;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.mockito.Mockito.*;

public class URLStencilClientWithCacheTest {

    private Map<String, DescriptorAndTypeName> descriptorMap;
    private static final String DESCRIPTOR_FILE_PATH = "__files/descriptors.bin";
    private static final String LOOKUP_KEY = "io.odpf.stencil.TestMessage";

    @Before
    public void setup() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        descriptorMap = new DescriptorMapBuilder().buildFrom(fileInputStream);
    }

    @Test
    public void getFromStencilClientSuccessfully() {
        DescriptorCacheLoader cacheLoader = mock(DescriptorCacheLoader.class);
        when(cacheLoader.load(LOOKUP_KEY)).thenReturn(descriptorMap);

        URLStencilClient stencilClient = new URLStencilClient(LOOKUP_KEY, ConfigFactory.create(StencilConfig.class, new HashMap<>()), cacheLoader);
        Descriptors.Descriptor result = stencilClient.get(LOOKUP_KEY);

        verify(cacheLoader, times(1)).load(LOOKUP_KEY);
        verify(cacheLoader, times(0)).reload(LOOKUP_KEY, descriptorMap);
        assertNotNull(result);
    }

    @Test(expected = StencilRuntimeException.class)
    public void getFromStencilClientOnException() throws Exception {
        DescriptorCacheLoader cacheLoader = mock(DescriptorCacheLoader.class);
        when(cacheLoader.load(LOOKUP_KEY)).thenThrow(new StencilRuntimeException(new Throwable()));

        URLStencilClient stencilClient = new URLStencilClient(LOOKUP_KEY, ConfigFactory.create(StencilConfig.class, new HashMap<>()), cacheLoader);
        stencilClient.get(LOOKUP_KEY);
    }

    @Test
    public void shouldNotRefreshCacheByDefault() {
        DescriptorCacheLoader cacheLoader = mock(DescriptorCacheLoader.class);
        when(cacheLoader.load(LOOKUP_KEY)).thenReturn(descriptorMap);

        FakeTicker fakeTicker = new FakeTicker();

        URLStencilClient stencilClient = new URLStencilClient(LOOKUP_KEY, ConfigFactory.create(StencilConfig.class, new HashMap<>()), cacheLoader, fakeTicker);
        Descriptors.Descriptor result = stencilClient.get(LOOKUP_KEY);

        fakeTicker.advance(1000, TimeUnit.MINUTES);
        Descriptors.Descriptor reloadedResult = stencilClient.get(LOOKUP_KEY);

        verify(cacheLoader, times(1)).load(LOOKUP_KEY);
        verify(cacheLoader, times(0)).reload(LOOKUP_KEY, descriptorMap);
        assertEquals(result, reloadedResult);
    }

    @Test
    public void getFromStencilClientSuccessfullySubsequentTimesFromCache() {
        DescriptorCacheLoader cacheLoader = mock(DescriptorCacheLoader.class);
        when(cacheLoader.load(LOOKUP_KEY)).thenReturn(descriptorMap);

        FakeTicker fakeTicker = new FakeTicker();
        Map<String, String> config = new HashMap<>();
        config.put("REFRESH_CACHE", "true");

        URLStencilClient stencilClient = new URLStencilClient(LOOKUP_KEY, ConfigFactory.create(StencilConfig.class, config), cacheLoader, fakeTicker);
        Descriptors.Descriptor result = stencilClient.get(LOOKUP_KEY);
        assertNotNull(result);

        fakeTicker.advance(stencilClient.getTTL().toMinutes() + 1, TimeUnit.MINUTES);
        Descriptors.Descriptor reloadedResult = stencilClient.get(LOOKUP_KEY);

        verify(cacheLoader, times(1)).load(LOOKUP_KEY);
        verify(cacheLoader, times(1)).reload(LOOKUP_KEY, descriptorMap);
        assertNotNull(reloadedResult);
    }

}
