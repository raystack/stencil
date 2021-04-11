package io.odpf.stencil.cache;

import com.google.common.io.ByteStreams;
import com.timgroup.statsd.NoOpStatsDClient;
import io.odpf.stencil.TestKey;
import io.odpf.stencil.exception.StencilRuntimeException;
import io.odpf.stencil.http.RemoteFile;
import io.odpf.stencil.models.DescriptorAndTypeName;
import org.apache.http.client.ClientProtocolException;
import org.junit.Test;

import java.io.FileInputStream;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

public class DescriptorCacheLoaderTest {
    private static final String DESCRIPTOR_FILE_PATH = "__files/descriptors.bin";
    private static final String LOOKUP_KEY = "io.odpf.stencil.TestMessage";
    private static final String TYPENAME_KEY = "io.odpf.stencil.TestMessage";


    @Test(expected = StencilRuntimeException.class)
    public void testStencilCacheLoadOnException() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        when(remoteFile.fetch(anyString())).thenThrow(new ClientProtocolException(""));
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(remoteFile, new NoOpStatsDClient(), null, true);
        cacheLoader.load(LOOKUP_KEY);
    }

    @Test
    public void testStencilCacheLoad() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        when(remoteFile.fetch(anyString())).thenReturn(bytes);

        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(remoteFile, new NoOpStatsDClient(), null, true);
        assertTrue(cacheLoader.load(LOOKUP_KEY).containsKey(LOOKUP_KEY));
    }

    @Test
    public void testStencilCacheReloadWithNewDescriptor() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        when(remoteFile.fetch(anyString())).thenReturn(bytes);

        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(remoteFile, new NoOpStatsDClient(), null, true);
        Map<String, DescriptorAndTypeName> prevDescriptor = new HashMap<>();
        assertTrue(cacheLoader.reload(LOOKUP_KEY, prevDescriptor).get().containsKey(LOOKUP_KEY));
    }

    @Test
    public void testStencilCacheReloadWithOldDescriptorOnException() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        when(remoteFile.fetch(anyString())).thenThrow(new ClientProtocolException(""));

        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(remoteFile, new NoOpStatsDClient(), null, true);

        Map<String, DescriptorAndTypeName> prevDescriptor = new HashMap<>();
        Map<String, DescriptorAndTypeName> result = cacheLoader.reload(LOOKUP_KEY, prevDescriptor).get();
        assertEquals(result.size(), 0);
    }

    @Test
    public void testStencilCacheReloadShouldCallOnProtoUpdateIfProtoChanges() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        ProtoUpdateListener protoUpdateListener = mock(ProtoUpdateListener.class);
        when(protoUpdateListener.getProto()).thenReturn(LOOKUP_KEY);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        when(remoteFile.fetch(LOOKUP_KEY)).thenReturn(bytes);
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(remoteFile, new NoOpStatsDClient(), protoUpdateListener, true);
        Map<String, DescriptorAndTypeName> prevDescriptor = new HashMap<>();
        prevDescriptor.put(LOOKUP_KEY, new DescriptorAndTypeName(TestKey.getDescriptor(), TYPENAME_KEY));
        assertTrue(cacheLoader.reload(LOOKUP_KEY, prevDescriptor).get().containsKey(LOOKUP_KEY));
        verify(protoUpdateListener, times(1)).onProtoUpdate(any(String.class), any(Map.class));
    }
}
