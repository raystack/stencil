package com.gojek.de.stencil.cache;

import com.gojek.de.stencil.exception.StencilRuntimeException;
import com.gojek.de.stencil.http.RemoteFile;
import com.gojek.stencil.TestMessage;
import com.google.common.io.ByteStreams;
import com.google.protobuf.Descriptors.Descriptor;
import com.timgroup.statsd.NoOpStatsDClient;
import org.apache.http.client.ClientProtocolException;
import org.junit.Test;

import java.io.FileInputStream;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

public class DescriptorCacheLoaderTest {
    private static final String DESCRIPTOR_FILE_PATH = "__files/descriptors.bin";
    private static final String LOOKUP_KEY = "com.gojek.stencil.TestMessage";

    @Test(expected = StencilRuntimeException.class)
    public void testStencilCacheLoadOnException() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        when(remoteFile.fetch(anyString())).thenThrow(new ClientProtocolException(""));
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(remoteFile, new NoOpStatsDClient(), null);
        cacheLoader.load(LOOKUP_KEY);
    }

    @Test
    public void testStencilCacheLoad() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        when(remoteFile.fetch(anyString())).thenReturn(bytes);

        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(remoteFile, new NoOpStatsDClient(), null);
        assertTrue(cacheLoader.load(LOOKUP_KEY).containsKey(LOOKUP_KEY));
    }

    @Test
    public void testStencilCacheReloadWithNewDescriptor() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        when(remoteFile.fetch(anyString())).thenReturn(bytes);

        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(remoteFile, new NoOpStatsDClient(), null);
        Map<String, Descriptor> prevDescriptor = new HashMap<>();
        assertTrue(cacheLoader.reload(LOOKUP_KEY, prevDescriptor).get().containsKey(LOOKUP_KEY));
    }

    @Test
    public void testStencilCacheReloadWithOldDescriptorOnException() throws Exception {
        RemoteFile remoteFile = mock(RemoteFile.class);
        when(remoteFile.fetch(anyString())).thenThrow(new ClientProtocolException(""));

        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(remoteFile, new NoOpStatsDClient(), null);

        Map<String, Descriptor> prevDescriptor = new HashMap<>();
        Map<String, Descriptor> result = cacheLoader.reload(LOOKUP_KEY, prevDescriptor).get();
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
        when(remoteFile.fetch(anyString())).thenReturn(bytes);
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(remoteFile, new NoOpStatsDClient(), protoUpdateListener);
        Map<String, Descriptor> prevDescriptor = new HashMap<>();
        TestMessage t = TestMessage.newBuilder().setSampleString("sample_value").build();
        prevDescriptor.put(LOOKUP_KEY, t.getDescriptorForType());
        assertTrue(cacheLoader.reload(LOOKUP_KEY, prevDescriptor).get().containsKey(LOOKUP_KEY));
        verify(protoUpdateListener, times(1)).onProtoUpdate();
    }
}
