package io.odpf.stencil.cache;

import com.google.common.io.ByteStreams;
import com.google.protobuf.Descriptors;
import com.timgroup.statsd.NoOpStatsDClient;
import io.odpf.stencil.SchemaUpdateListener;
import io.odpf.stencil.TestKey;
import io.odpf.stencil.exception.StencilRuntimeException;
import io.odpf.stencil.http.RemoteFile;
import org.apache.http.client.ClientProtocolException;
import org.junit.After;
import org.junit.Test;
import org.mockito.Mock;
import java.io.FileInputStream;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

public class DescriptorCacheLoaderTest {
    private SchemaCacheLoader cacheLoader;
    private static final String DESCRIPTOR_FILE_PATH = "__files/descriptors.bin";
    private static final String LOOKUP_KEY = "io.odpf.stencil.TestMessage";

    @After
    public void teatDown() throws Exception {
        cacheLoader.close();
    }


    @Test(expected = StencilRuntimeException.class)
    public void testStencilCacheLoadOnException() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        when(remoteFile.fetch(anyString())).thenThrow(new ClientProtocolException(""));
        cacheLoader = new SchemaCacheLoader(remoteFile, new NoOpStatsDClient(), null, SchemaRefreshStrategy.longPollingStrategy(), true);
        cacheLoader.load(LOOKUP_KEY);
    }

    @Test
    public void testStencilCacheLoad() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        when(remoteFile.fetch(anyString())).thenReturn(bytes);

        cacheLoader = new SchemaCacheLoader(remoteFile, new NoOpStatsDClient(), null, SchemaRefreshStrategy.longPollingStrategy(), true);
        assertTrue(cacheLoader.load(LOOKUP_KEY).containsKey(LOOKUP_KEY));
    }

    @Test
    public void testStencilCacheReloadWithNewDescriptor() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        when(remoteFile.fetch(anyString())).thenReturn(bytes);

        cacheLoader = new SchemaCacheLoader(remoteFile, new NoOpStatsDClient(), null, SchemaRefreshStrategy.longPollingStrategy(), true);
        Map<String, Descriptors.Descriptor> prevDescriptor = new HashMap<>();
        assertTrue(cacheLoader.reload(LOOKUP_KEY, prevDescriptor).get().containsKey(LOOKUP_KEY));
    }

    @Test
    public void testStencilCacheReloadWithOldDescriptorOnException() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        when(remoteFile.fetch(anyString())).thenThrow(new ClientProtocolException(""));

        cacheLoader = new SchemaCacheLoader(remoteFile, new NoOpStatsDClient(), null, SchemaRefreshStrategy.longPollingStrategy(), true);

        Map<String, Descriptors.Descriptor> prevDescriptor = new HashMap<>();
        Map<String, Descriptors.Descriptor> result = cacheLoader.reload(LOOKUP_KEY, prevDescriptor).get();
        assertEquals(result.size(), 0);
    }

    @Test
    public void testStencilCacheReloadShouldCallOnProtoUpdateIfProtoChanges() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        SchemaUpdateListener mockedListener = mock(SchemaUpdateListener.class);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        when(remoteFile.fetch(LOOKUP_KEY)).thenReturn(bytes);
        cacheLoader = new SchemaCacheLoader(remoteFile, new NoOpStatsDClient(), mockedListener, SchemaRefreshStrategy.longPollingStrategy(), true);
        Map<String, Descriptors.Descriptor> prevDescriptor = new HashMap<>();
        prevDescriptor.put(LOOKUP_KEY, TestKey.getDescriptor());
        assertTrue(cacheLoader.reload(LOOKUP_KEY, prevDescriptor).get().containsKey(LOOKUP_KEY));
        verify(mockedListener).onSchemaUpdate(any(Map.class));
    }
}
