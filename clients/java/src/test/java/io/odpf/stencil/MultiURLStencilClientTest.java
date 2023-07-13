package org.raystack.stencil;

import com.github.tomakehurst.wiremock.WireMockServer;
import com.github.tomakehurst.wiremock.core.WireMockConfiguration;
import com.google.protobuf.Descriptors;
import org.raystack.stencil.client.StencilClient;
import org.raystack.stencil.config.StencilConfig;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import java.util.ArrayList;
import java.util.Map;

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
        StencilClient c = StencilClientFactory.getClient(urls, StencilConfig.builder().build());
        Map<String, Descriptors.Descriptor> descMap = c.getAll();
        assertNotNull(descMap);
        Descriptors.Descriptor desc = c.get("org.raystack.stencil.TestMessage");
        assertNotNull(desc);
    }
}
