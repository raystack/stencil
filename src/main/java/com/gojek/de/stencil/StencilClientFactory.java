package com.gojek.de.stencil;

import com.gojek.de.stencil.cache.DescriptorCacheLoader;
import com.google.common.cache.CacheLoader;

import java.util.List;
import java.util.Map;

public class StencilClientFactory {
    public static StencilClient getClient(String url, Map<String, String> config) {

        CacheLoader cacheLoader = new DescriptorCacheLoader(config, new RemoteFileImpl());
        return new URLStencilClient(url, config, cacheLoader);
    }

    public static StencilClient getClient(List<String> urls, Map<String, String> config) {
        CacheLoader cacheLoader = new DescriptorCacheLoader(config, new RemoteFileImpl());
        return new MultiURLStencilClient(urls, config, cacheLoader);
    }

    public static StencilClient getClient() {
        return new ClassLoadStencilClient();
    }
}
