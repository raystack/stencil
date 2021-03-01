package com.gojek.de.stencil;

import com.gojek.de.stencil.cache.DescriptorCacheLoader;
import com.gojek.de.stencil.cache.ProtoUpdateListener;
import com.gojek.de.stencil.client.ClassLoadStencilClient;
import com.gojek.de.stencil.client.MultiURLStencilClient;
import com.gojek.de.stencil.client.StencilClient;
import com.gojek.de.stencil.client.URLStencilClient;
import com.gojek.de.stencil.config.StencilConfig;
import com.gojek.de.stencil.http.RemoteFileImpl;
import com.gojek.de.stencil.http.RetryHttpClient;
import com.timgroup.statsd.NoOpStatsDClient;
import com.timgroup.statsd.StatsDClient;
import org.aeonbits.owner.ConfigFactory;

import java.util.List;
import java.util.Map;

public class StencilClientFactory {
    public static StencilClient getClient(String url, Map<String, String> config, StatsDClient statsDClient) {
        StencilConfig stencilConfig = ConfigFactory.create(StencilConfig.class, config);
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, null, stencilConfig.shouldAutoRefreshCache());
        return new URLStencilClient(url, stencilConfig, cacheLoader);
    }

    public static StencilClient getClient(String url, Map<String, String> config, StatsDClient statsDClient, ProtoUpdateListener protoUpdateListener) {
        StencilConfig stencilConfig = ConfigFactory.create(StencilConfig.class, config);
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, protoUpdateListener, stencilConfig.shouldAutoRefreshCache());
        return new URLStencilClient(url, stencilConfig, cacheLoader);
    }

    public static StencilClient getClient(List<String> urls, Map<String, String> config, StatsDClient statsDClient) {
        StencilConfig stencilConfig = ConfigFactory.create(StencilConfig.class, config);
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, null, stencilConfig.shouldAutoRefreshCache());
        return new MultiURLStencilClient(urls, stencilConfig, cacheLoader);
    }

    public static StencilClient getClient(List<String> urls, Map<String, String> config, StatsDClient statsDClient, ProtoUpdateListener protoUpdateListener) {
        StencilConfig stencilConfig = ConfigFactory.create(StencilConfig.class, config);
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, protoUpdateListener, stencilConfig.shouldAutoRefreshCache());
        return new MultiURLStencilClient(urls, stencilConfig, cacheLoader);
    }

    public static StencilClient getClient(String url, Map<String, String> config) {
        return getClient(url, config, new NoOpStatsDClient());
    }

    public static StencilClient getClient(List<String> urls, Map<String, String> config) {
        return getClient(urls, config, new NoOpStatsDClient());
    }

    public static StencilClient getClient() {
        return new ClassLoadStencilClient();
    }
}
