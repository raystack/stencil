package com.gotocompany.stencil.http;

import com.github.tomakehurst.wiremock.junit.WireMockRule;
import com.gotocompany.stencil.config.StencilConfig;

import org.apache.http.Header;
import org.apache.http.HttpHeaders;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.message.BasicHeader;
import org.junit.Rule;
import org.junit.Test;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import static com.github.tomakehurst.wiremock.client.WireMock.*;


public class RetryHttpClientTest {

    @Rule
    public WireMockRule service = new WireMockRule(8081);

    @Test
    public void shouldUseAuthenticationBearerTokenFromStencilConfig() throws IOException {
        String token = "test-token";

        service.stubFor(any(anyUrl())
                .willReturn(aResponse()
                .withStatus(200))
        );
        Header authHeader = new BasicHeader(HttpHeaders.AUTHORIZATION, "Bearer " + token);
        List<Header> headers = new ArrayList<Header>();
        headers.add(authHeader);

        CloseableHttpClient httpClient = new RetryHttpClient().create(StencilConfig.builder().fetchHeaders(headers).build());
        httpClient.execute(new HttpGet(service.url("/test/stencil/auth/header")));

        verify(getRequestedFor(anyUrl()).withHeader(HttpHeaders.AUTHORIZATION, equalTo("Bearer " + token)));
    }
}
