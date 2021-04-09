package io.odpf.stencil;

import io.odpf.stencil.cache.DescriptorCacheLoader;
import io.odpf.stencil.cache.ProtoUpdateListener;
import io.odpf.stencil.client.ClassLoadStencilClient;
import io.odpf.stencil.client.MultiURLStencilClient;
import io.odpf.stencil.client.StencilClient;
import io.odpf.stencil.client.URLStencilClient;
import io.odpf.stencil.config.StencilConfig;
import io.odpf.stencil.http.RemoteFileImpl;
import io.odpf.stencil.http.RetryHttpClient;
import com.timgroup.statsd.NoOpStatsDClient;
import com.timgroup.statsd.StatsDClient;
import org.aeonbits.owner.ConfigFactory;

import java.util.List;
import java.util.Map;

public class StencilClientFactory {
    public static StencilClient getClient(String url, StencilConfig config, StatsDClient statsDClient) {
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, null, config.getCacheAutoRefresh());
        return new URLStencilClient(url, config, cacheLoader);
    }

    public static StencilClient getClient(String url, StencilConfig config, StatsDClient statsDClient, ProtoUpdateListener protoUpdateListener) {
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, protoUpdateListener, config.getCacheAutoRefresh());
        return new URLStencilClient(url, config, cacheLoader);
    }

    public static StencilClient getClient(List<String> urls, StencilConfig config, StatsDClient statsDClient) {
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, null, config.getCacheAutoRefresh());
        return new MultiURLStencilClient(urls, config, cacheLoader);
    }

    public static StencilClient getClient(List<String> urls, StencilConfig config, StatsDClient statsDClient, ProtoUpdateListener protoUpdateListener) {
        DescriptorCacheLoader cacheLoader = new DescriptorCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, protoUpdateListener, config.getCacheAutoRefresh());
        return new MultiURLStencilClient(urls, config, cacheLoader);
    }

    public static StencilClient getClient(String url, StencilConfig config) {
        return getClient(url, config, new NoOpStatsDClient());
    }

    public static StencilClient getClient(List<String> urls, StencilConfig config) {
        return getClient(urls, config, new NoOpStatsDClient());
    }

    public static StencilClient getClient() {
        return new ClassLoadStencilClient();
    }
}
