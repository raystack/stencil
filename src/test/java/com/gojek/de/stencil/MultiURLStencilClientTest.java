package com.gojek.de.stencil;

import com.github.tomakehurst.wiremock.WireMockServer;
import com.github.tomakehurst.wiremock.core.WireMockConfiguration;
import com.gojek.de.stencil.client.StencilClient;
import com.google.protobuf.Descriptors;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

import static org.junit.Assert.assertNotNull;

public class MultiURLStencilClientTest {
    private WireMockServer wireMockServer;

    @Before
    public void setup() {
        WireMockConfiguration config = new WireMockConfiguration();
        config = config.withRootDirectory("src/test/resources/").port(8082);
        wireMockServer = new WireMockServer(config);
        wireMockServer.start();
    }

    @After
    public void tearDown() {
        wireMockServer.stop();
    }

    @Test
    public void shouldReturnDescriptor() {
        ArrayList<String> urls = new ArrayList<String>();
        urls.add("http://localhost:8082/descriptors.bin");
        Map<String, String> config = new HashMap<>();
        StencilClient c = StencilClientFactory.getClient(urls, config);
        Descriptors.Descriptor desc = c.get("com.gojek.stencil.TestMessage");
        assertNotNull(desc);
    }
}
