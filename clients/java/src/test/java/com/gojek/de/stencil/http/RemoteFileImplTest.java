package com.gojek.de.stencil.http;

import org.apache.http.HttpEntity;
import org.apache.http.StatusLine;
import org.apache.http.client.ClientProtocolException;
import org.apache.http.client.ResponseHandler;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpUriRequest;
import org.apache.http.impl.client.CloseableHttpClient;
import org.junit.Test;

import java.io.ByteArrayInputStream;
import java.io.IOException;

import static org.junit.Assert.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

public class RemoteFileImplTest {
    @Test
    public void testRemoteFileFetchOnSuccess() throws Exception {
        CloseableHttpClient httpClient = mock(CloseableHttpClient.class);
        when(httpClient.execute(any(HttpUriRequest.class), any(ResponseHandler.class))).thenReturn("test".getBytes());

        RemoteFileImpl remoteFileImpl = new RemoteFileImpl(httpClient);
        String result = new String(remoteFileImpl.fetch("url"));
        assertEquals(result, "test");
    }

    @Test(expected = ClientProtocolException.class)
    public void testRemoteFileFetchOnException() throws Exception {
        CloseableHttpClient httpClient = mock(CloseableHttpClient.class);

        when(httpClient.execute(any(HttpUriRequest.class), any(ResponseHandler.class))).thenThrow(new ClientProtocolException(""));
        new RemoteFileImpl(httpClient).fetch("");
    }

    @Test
    public void testHandleResponseWhen200() throws IOException {
        CloseableHttpClient httpClient = mock(CloseableHttpClient.class);
        CloseableHttpResponse httpResponse = mock(CloseableHttpResponse.class);
        StatusLine statusLine = mock(StatusLine.class);
        HttpEntity entity = mock(HttpEntity.class);

        when(statusLine.getStatusCode()).thenReturn(200);
        when(httpResponse.getStatusLine()).thenReturn(statusLine);
        when(entity.getContent()).thenReturn(new ByteArrayInputStream("test".getBytes()));
        when(httpResponse.getEntity()).thenReturn(entity);

        RemoteFileImpl remoteFileImpl = new RemoteFileImpl(httpClient);
        assertEquals(new String(remoteFileImpl.handleResponse(httpResponse)), "test");
    }

    @Test(expected = ClientProtocolException.class)
    public void testHandleResponseWhen500() throws IOException {
        CloseableHttpClient httpClient = mock(CloseableHttpClient.class);
        CloseableHttpResponse httpResponse = mock(CloseableHttpResponse.class);
        StatusLine statusLine = mock(StatusLine.class);
        HttpEntity entity = mock(HttpEntity.class);

        when(statusLine.getStatusCode()).thenReturn(503);
        when(httpResponse.getStatusLine()).thenReturn(statusLine);
        when(entity.getContent()).thenReturn(new ByteArrayInputStream("test".getBytes()));
        when(httpResponse.getEntity()).thenReturn(entity);

        RemoteFileImpl remoteFileImpl = new RemoteFileImpl(httpClient);
        assertEquals(new String(remoteFileImpl.handleResponse(httpResponse)), "test");
    }
}
