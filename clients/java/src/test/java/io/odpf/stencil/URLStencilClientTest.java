package io.odpf.stencil;

import com.github.tomakehurst.wiremock.WireMockServer;
import com.github.tomakehurst.wiremock.core.WireMockConfiguration;
import com.google.protobuf.Descriptors;
import io.odpf.stencil.client.StencilClient;
import io.odpf.stencil.config.StencilConfig;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import java.io.IOException;
import java.util.Map;

import static org.junit.Assert.assertNotNull;

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
        StencilClient c = StencilClientFactory.getClient(url, StencilConfig.builder().build());
        Map<String, Descriptors.Descriptor> descMap = c.getAll();
        assertNotNull(descMap);
        Descriptors.Descriptor desc = c.get("io.odpf.stencil.TestMessage");
        assertNotNull(desc);
        c.refresh();
        c.close();
    }

}
