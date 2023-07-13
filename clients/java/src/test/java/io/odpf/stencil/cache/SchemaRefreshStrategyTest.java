package org.raystack.stencil.cache;

import static org.junit.Assert.assertNotSame;
import static org.junit.Assert.assertSame;
import static org.junit.Assert.assertTrue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;
import com.google.common.io.ByteStreams;
import com.google.protobuf.Descriptors;

import org.junit.Test;

import org.raystack.stencil.http.RemoteFile;

public class SchemaRefreshStrategyTest {
    private static final String TYPENAME_KEY = "org.raystack.stencil.TestMessage";
    private static final String DESCRIPTOR_FILE_PATH = "__files/descriptors.bin";
    private static final String baseURL = "http://localhost:8000/v1beta1/namespaces/protobuf-test/schemas/esb-log-entities";
    private static final String versionsURL = String.format("%s/versions", baseURL);

    @Test
    public void testLongPollingStrategyShouldReturnPreviousDataIfVersionsEmpty() throws IOException {
        byte[] versionsData = "{\"versions\": []}".toString().getBytes("utf-8");
        RemoteFile remoteFile = mock(RemoteFile.class);
        when(remoteFile.fetch(versionsURL)).thenReturn(versionsData);
        Map<String, Descriptors.Descriptor> prevData = new HashMap<>();
        SchemaRefreshStrategy fn = SchemaRefreshStrategy.versionBasedRefresh();
        Map<String, Descriptors.Descriptor> newData = fn.refresh(baseURL, remoteFile, prevData);
        assertSame(prevData, newData);
        verify(remoteFile).fetch(versionsURL);
    }

    @Test
    public void testLongPollingStrategyShouldGetMaxVersionData() throws IOException {
        String versionedURL = String.format("%s/%d", versionsURL, 1);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        byte[] versionsData = "{\"versions\": [1]}".toString().getBytes("utf-8");
        RemoteFile remoteFile = mock(RemoteFile.class);
        when(remoteFile.fetch(versionsURL)).thenReturn(versionsData);
        when(remoteFile.fetch(versionedURL)).thenReturn(bytes);
        Map<String, Descriptors.Descriptor> prevData = new HashMap<>();
        SchemaRefreshStrategy fn = SchemaRefreshStrategy.versionBasedRefresh();
        Map<String, Descriptors.Descriptor> newData = fn.refresh(baseURL, remoteFile, prevData);
        verify(remoteFile).fetch(versionsURL);
        verify(remoteFile).fetch(versionedURL);
        assertTrue(newData.containsKey(TYPENAME_KEY));
    }

    @Test
    public void testLongPollingStrategyShouldRememberOldVersion() throws IOException {
        String versionedURL = String.format("%s/%d", versionsURL, 1);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        byte[] versionsData = "{\"versions\": [1]}".toString().getBytes("utf-8");
        RemoteFile remoteFile = mock(RemoteFile.class);
        when(remoteFile.fetch(versionsURL)).thenReturn(versionsData);
        when(remoteFile.fetch(versionedURL)).thenReturn(bytes);
        Map<String, Descriptors.Descriptor> prevData = new HashMap<>();
        SchemaRefreshStrategy fn = SchemaRefreshStrategy.versionBasedRefresh();
        Map<String, Descriptors.Descriptor> newData = fn.refresh(baseURL, remoteFile, prevData);
        Map<String, Descriptors.Descriptor> updated = fn.refresh(baseURL, remoteFile, newData);
        verify(remoteFile, times(2)).fetch(versionsURL);
        verify(remoteFile).fetch(versionedURL);
        assertSame(newData, updated);
        assertTrue(newData.containsKey(TYPENAME_KEY));
    }

    @Test
    public void testLongPollingStrategyShouldUpdateDescriptorDataOnlyIfVersionChanged() throws IOException {
        String versionedURL = String.format("%s/%d", versionsURL, 1);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        byte[] versionsData = "{\"versions\": [1]}".toString().getBytes("utf-8");
        byte[] updatedVersions = "{\"versions\": [1, 2]}".toString().getBytes("utf-8");
        RemoteFile remoteFile = mock(RemoteFile.class);
        when(remoteFile.fetch(versionsURL)).thenReturn(versionsData).thenReturn(updatedVersions);
        when(remoteFile.fetch(versionedURL)).thenReturn(bytes);
        when(remoteFile.fetch(String.format("%s/%d", versionsURL, 2))).thenReturn(bytes);
        Map<String, Descriptors.Descriptor> prevData = new HashMap<>();
        SchemaRefreshStrategy fn = SchemaRefreshStrategy.versionBasedRefresh();
        Map<String, Descriptors.Descriptor> newData = fn.refresh(baseURL, remoteFile, prevData);
        Map<String, Descriptors.Descriptor> updated = fn.refresh(baseURL, remoteFile, newData);
        verify(remoteFile, times(2)).fetch(versionsURL);
        verify(remoteFile).fetch(versionedURL);
        verify(remoteFile).fetch(String.format("%s/%d", versionsURL, 2));
        assertNotSame(newData, updated);
        assertTrue(newData.containsKey(TYPENAME_KEY));
    }

    @Test
    public void testLongPollingStrategyShouldNotShareDataBetweenDifferentInstances() throws IOException {
        String versionedURL = String.format("%s/%d", versionsURL, 1);
        ClassLoader classLoader = getClass().getClassLoader();
        InputStream fileInputStream = new FileInputStream(classLoader.getResource(DESCRIPTOR_FILE_PATH).getFile());
        byte[] bytes = ByteStreams.toByteArray(fileInputStream);
        byte[] versionsData = "{\"versions\": [1]}".toString().getBytes("utf-8");
        byte[] updatedVersions = "{\"versions\": [1, 2]}".toString().getBytes("utf-8");
        RemoteFile remoteFile = mock(RemoteFile.class);
        when(remoteFile.fetch(versionsURL)).thenReturn(versionsData, updatedVersions, updatedVersions);
        when(remoteFile.fetch(versionedURL)).thenReturn(bytes);
        when(remoteFile.fetch(String.format("%s/%d", versionsURL, 2))).thenReturn(bytes);
        Map<String, Descriptors.Descriptor> prevData = new HashMap<>();
        SchemaRefreshStrategy fn = SchemaRefreshStrategy.versionBasedRefresh();
        SchemaRefreshStrategy fn2 = SchemaRefreshStrategy.versionBasedRefresh();
        Map<String, Descriptors.Descriptor> newData = fn.refresh(baseURL, remoteFile, prevData);
        Map<String, Descriptors.Descriptor> updated = fn2.refresh(baseURL, remoteFile, newData);
        Map<String, Descriptors.Descriptor> updatedData = fn2.refresh(baseURL, remoteFile, updated);
        verify(remoteFile, times(3)).fetch(versionsURL);
        verify(remoteFile).fetch(versionedURL);
        verify(remoteFile).fetch(String.format("%s/%d", versionsURL, 2));
        assertNotSame(newData, updated);
        assertNotSame(prevData, updated);
        assertSame(updated, updatedData);
    }
}
