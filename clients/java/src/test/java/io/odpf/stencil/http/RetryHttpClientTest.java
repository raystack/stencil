package io.odpf.stencil.http;

import com.github.tomakehurst.wiremock.junit.WireMockRule;
import io.odpf.stencil.config.StencilConfig;
import org.apache.http.HttpHeaders;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.junit.Rule;
import org.junit.Test;

import java.io.IOException;

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

        CloseableHttpClient httpClient = new RetryHttpClient().create(StencilConfig.builder().fetchAuthBearerToken(token).build());
        httpClient.execute(new HttpGet(service.url("/test/stencil/auth/header")));

        verify(getRequestedFor(anyUrl()).withHeader(HttpHeaders.AUTHORIZATION, equalTo("Bearer " + token)));
    }
}
