package com.gojek.de.stencil;

import com.github.tomakehurst.wiremock.WireMockServer;
import com.github.tomakehurst.wiremock.core.WireMockConfiguration;
import com.gojek.de.stencil.client.StencilClient;
import com.google.protobuf.Descriptors;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

import static org.junit.Assert.*;

public class URLStencilClientTest {
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
    public void downloadFile() throws IOException {
        String url = "http://localhost:8082/descriptors.bin";
        Map<String, String> config = new HashMap<>();
        StencilClient c = StencilClientFactory.getClient(url, config);
        Map<String, Descriptors.Descriptor> descMap = c.getAll();
        assertNotNull(descMap);
        Descriptors.Descriptor desc = c.get("com.gojek.stencil.TestMessage");
        assertNotNull(desc);
        c.refresh();
        c.close();
    }

}
