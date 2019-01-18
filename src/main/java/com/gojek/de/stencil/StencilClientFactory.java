package com.gojek.de.stencil;

import com.gojek.de.stencil.cache.DescriptorCacheLoader;
import com.gojek.de.stencil.client.ClassLoadStencilClient;
import com.gojek.de.stencil.client.MultiURLStencilClient;
import com.gojek.de.stencil.client.StencilClient;
import com.gojek.de.stencil.client.URLStencilClient;
import com.gojek.de.stencil.http.RemoteFileImpl;
import com.gojek.de.stencil.http.RetryHttpClient;
import com.timgroup.statsd.StatsDClient;

import java.util.List;
import java.util.Map;
import java.util.Optional;

public class StencilClientFactory {
    public static StencilClient getClient(String url, Map<String, String> config, Optional<StatsDClient> statsDClientOpt) {
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClientOpt);
        return new URLStencilClient(url, config, cacheLoader);
    }

    public static StencilClient getClient(List<String> urls, Map<String, String> config, Optional<StatsDClient> statsDClientOpt) {
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClientOpt);
        return new MultiURLStencilClient(urls, config, cacheLoader);
    }

    public static StencilClient getClient(String url, Map<String, String> config) {
        return getClient(url, config, Optional.empty());
    }

    public static StencilClient getClient(List<String> urls, Map<String, String> config) {
        return getClient(urls, config, Optional.empty());
    }

    public static StencilClient getClient() {
        return new ClassLoadStencilClient();
    }
}
