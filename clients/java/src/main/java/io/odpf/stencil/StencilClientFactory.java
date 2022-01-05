package io.odpf.stencil;

import com.timgroup.statsd.NoOpStatsDClient;
import com.timgroup.statsd.StatsDClient;
import io.odpf.stencil.cache.SchemaCacheLoader;
import io.odpf.stencil.cache.ProtoUpdateListener;
import io.odpf.stencil.client.ClassLoadStencilClient;
import io.odpf.stencil.client.MultiURLStencilClient;
import io.odpf.stencil.client.StencilClient;
import io.odpf.stencil.client.URLStencilClient;
import io.odpf.stencil.config.StencilConfig;
import io.odpf.stencil.http.RemoteFileImpl;
import io.odpf.stencil.http.RetryHttpClient;

import java.util.List;


/**
 * Provides static methods for the creation of {@link io.odpf.stencil.client.StencilClient}
 * object with configurations and various options like
 * single URLs, multiple URLs, statsd client for monitoring
 * and {@link io.odpf.stencil.cache.ProtoUpdateListener} for callback on proto schema update.
 */
public class StencilClientFactory {
    public static StencilClient getClient(String url, StencilConfig config, StatsDClient statsDClient) {
        SchemaCacheLoader cacheLoader = new SchemaCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, null, config.getRefreshStrategy(), config.getCacheAutoRefresh());
        return new URLStencilClient(url, config, cacheLoader);
    }

    /**
     * @param url URL to fetch and cache protobuf descriptor set in the client
     * @param config Stencil configs
     * @param statsDClient StatsD client to push metrics
     * @param protoUpdateListener {@link io.odpf.stencil.cache.ProtoUpdateListener#onProtoUpdate}
     *                            will be called when schema gets updated
     * @return Stencil client for single URL, statsd client
     * for monitoring and ProtoUpdateListener for callback
     */
    public static StencilClient getClient(String url, StencilConfig config, StatsDClient statsDClient, ProtoUpdateListener protoUpdateListener) {
        SchemaCacheLoader cacheLoader = new SchemaCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, protoUpdateListener, config.getRefreshStrategy(), config.getCacheAutoRefresh());
        return new URLStencilClient(url, config, cacheLoader);
    }

    /**
     * @param urls List of URLs to fetch and cache protobuf descriptor sets in the client
     * @param config Stencil configs
     * @param statsDClient StatsD client to push metrics
     * @return Stencil client for multiple URLs and statsd client for monitoring
     */
    public static StencilClient getClient(List<String> urls, StencilConfig config, StatsDClient statsDClient) {
        SchemaCacheLoader cacheLoader = new SchemaCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, null, config.getRefreshStrategy(), config.getCacheAutoRefresh());
        return new MultiURLStencilClient(urls, config, cacheLoader);
    }

    /**
     * @param urls List of URLs to fetch and cache protobuf descriptor sets in the client
     * @param config Stencil configs
     * @param statsDClient StatsD client to push metrics
     * @param protoUpdateListener {@link io.odpf.stencil.cache.ProtoUpdateListener#onProtoUpdate}
     *                            will be called when schema gets updated
     * @return Stencil client for multiple URLs, statsd client
     * for monitoring and ProtoUpdateListener for callback
     */
    public static StencilClient getClient(List<String> urls, StencilConfig config, StatsDClient statsDClient, ProtoUpdateListener protoUpdateListener) {
        SchemaCacheLoader cacheLoader = new SchemaCacheLoader(new RemoteFileImpl(new RetryHttpClient().create(config)), statsDClient, protoUpdateListener, config.getRefreshStrategy(), config.getCacheAutoRefresh());
        return new MultiURLStencilClient(urls, config, cacheLoader);
    }

    /**
     * @param url URL to fetch and cache protobuf descriptor set in the client
     * @param config Stencil configs
     * @return Stencil client for single URL
     */
    public static StencilClient getClient(String url, StencilConfig config) {
        return getClient(url, config, new NoOpStatsDClient());
    }

    /**
     * @param urls List of URLs to fetch and cache protobuf descriptor sets in the client
     * @param config Stencil configs
     * @return Stencil client for multiple URLs
     */
    public static StencilClient getClient(List<String> urls, StencilConfig config) {
        return getClient(urls, config, new NoOpStatsDClient());
    }

    /**
     * @return Stencil client for getting descriptors from classes in classpath
     */
    public static StencilClient getClient() {
        return new ClassLoadStencilClient();
    }
}
