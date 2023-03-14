package com.gotocompany.stencil.http;

import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.client.ClientProtocolException;
import org.apache.http.client.ResponseHandler;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.util.EntityUtils;

import java.io.IOException;

public class RemoteFileImpl implements RemoteFile, ResponseHandler<byte[]> {
    private CloseableHttpClient closeableHttpClient;

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
    public byte[] handleResponse(HttpResponse response) throws IOException {
        int status = response.getStatusLine().getStatusCode();
        if (status >= 200 && status < 300) {
            HttpEntity entity = response.getEntity();
            return entity != null ? EntityUtils.toByteArray(entity) : null;
        } else {
            throw new ClientProtocolException("Unexpected response status: " + status);
        }
    }
}
