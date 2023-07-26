package com.gotocompany.stencil.http;

import org.apache.hc.client5.http.ClientProtocolException;
import org.apache.hc.client5.http.classic.methods.HttpGet;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.core5.http.ClassicHttpResponse;
import org.apache.hc.core5.http.HttpEntity;
import org.apache.hc.core5.http.io.HttpClientResponseHandler;
import org.apache.hc.core5.http.io.entity.EntityUtils;

import java.io.IOException;

public class RemoteFileImpl implements RemoteFile, HttpClientResponseHandler<byte[]> {
    private final CloseableHttpClient closeableHttpClient;

    public RemoteFileImpl(CloseableHttpClient httpClient) {
        this.closeableHttpClient = httpClient;
    }

    public byte[] fetch(String url) throws IOException {
        HttpGet httpget = new HttpGet(url);
        byte[] responseBody;
        responseBody = closeableHttpClient.execute(httpget, this);
        return responseBody;
    }

    @Override
    public void close() throws IOException {
        closeableHttpClient.close();
    }

    @Override
    public byte[] handleResponse(ClassicHttpResponse response) throws IOException {
        int status = response.getCode();
        if (status >= 200 && status < 300) {
            HttpEntity entity = response.getEntity();
            return entity != null ? EntityUtils.toByteArray(entity) : null;
        } else {
            throw new ClientProtocolException("Unexpected response status: " + status);
        }
    }
}
