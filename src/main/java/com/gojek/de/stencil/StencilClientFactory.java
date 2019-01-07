package com.gojek.de.stencil;

import java.util.List;
import java.util.Map;

public class StencilClientFactory {
    public static StencilClient getClient(String url, Map<String, String> config) {
        return new URLStencilClient(url, config);
    }

    public static StencilClient getClient(List<String> urls, Map<String, String> config) {
        return new MultiURLStencilClient(urls, config);
    }

    public static StencilClient getClient() {
        return new ClassLoadStencilClient();
    }
}
