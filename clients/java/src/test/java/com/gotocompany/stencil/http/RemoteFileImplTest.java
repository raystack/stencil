package com.gotocompany.stencil.http;

import org.apache.hc.client5.http.ClientProtocolException;
import org.apache.hc.client5.http.classic.methods.HttpUriRequest;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.core5.http.ClassicHttpResponse;
import org.apache.hc.core5.http.HttpEntity;
import org.apache.hc.core5.http.HttpException;
import org.apache.hc.core5.http.io.HttpClientResponseHandler;
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
        when(httpClient.execute(any(HttpUriRequest.class), any(HttpClientResponseHandler.class))).thenReturn("test".getBytes());

        RemoteFileImpl remoteFileImpl = new RemoteFileImpl(httpClient);
        String result = new String(remoteFileImpl.fetch("url"));
        assertEquals(result, "test");
    }

    @Test(expected = ClientProtocolException.class)
    public void testRemoteFileFetchOnException() throws Exception {
        CloseableHttpClient httpClient = mock(CloseableHttpClient.class);

        when(httpClient.execute(any(HttpUriRequest.class), any(HttpClientResponseHandler.class))).thenThrow(new ClientProtocolException(""));
        new RemoteFileImpl(httpClient).fetch("");
    }

    @Test
    public void testHandleResponseWhen200() throws IOException, HttpException {
        CloseableHttpClient httpClient = mock(CloseableHttpClient.class);
        ClassicHttpResponse httpResponse = mock(ClassicHttpResponse.class);
        HttpEntity entity = mock(HttpEntity.class);

        when(httpResponse.getCode()).thenReturn(200);
        when(entity.getContent()).thenReturn(new ByteArrayInputStream("test".getBytes()));
        when(httpResponse.getEntity()).thenReturn(entity);

        RemoteFileImpl remoteFileImpl = new RemoteFileImpl(httpClient);
        assertEquals(new String(remoteFileImpl.handleResponse(httpResponse)), "test");
    }

    @Test(expected = ClientProtocolException.class)
    public void testHandleResponseWhen500() throws IOException, HttpException {
        CloseableHttpClient httpClient = mock(CloseableHttpClient.class);
        ClassicHttpResponse httpResponse = mock(ClassicHttpResponse.class);
        HttpEntity entity = mock(HttpEntity.class);

        when(httpResponse.getCode()).thenReturn(503);
        when(entity.getContent()).thenReturn(new ByteArrayInputStream("test".getBytes()));
        when(httpResponse.getEntity()).thenReturn(entity);

        RemoteFileImpl remoteFileImpl = new RemoteFileImpl(httpClient);
        assertEquals(new String(remoteFileImpl.handleResponse(httpResponse)), "test");
    }
}
