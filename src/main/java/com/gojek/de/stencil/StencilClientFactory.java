package com.gojek.de.stencil;

import java.util.Map;

public class StencilClientFactory {
    public static StencilClient getClient(String url, Map<String, String> options) {
        return new URLStencilClient(url, options);
    }

    public static StencilClient getClient() {
        return new ClassLoadStencilClient();
    }
}
