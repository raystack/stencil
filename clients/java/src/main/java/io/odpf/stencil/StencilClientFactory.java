package io.odpf.stencil;

import io.odpf.stencil.cache.SchemaCacheLoader;
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
 * and {@link io.odpf.stencil.SchemaUpdateListener} for callback on schema change.
 */
public class StencilClientFactory {
    /**
     * @param url URL to fetch and cache protobuf descriptor set in the client
     * @param config Stencil configs
     * @return Stencil client for single URL
     */
    public static StencilClient getClient(String url, StencilConfig config) {
        SchemaCacheLoader cacheLoader = new SchemaCacheLoader(new RemoteFileImpl(RetryHttpClient.create(config)), config);
        return new URLStencilClient(url, config, cacheLoader);
    }

    /**
     * @param urls List of URLs to fetch and cache protobuf descriptor sets in the client
     * @param config Stencil configs
     * @return Stencil client for multiple URLs
     */
    public static StencilClient getClient(List<String> urls, StencilConfig config) {
        SchemaCacheLoader cacheLoader = new SchemaCacheLoader(new RemoteFileImpl(RetryHttpClient.create(config)), config);
        return new MultiURLStencilClient(urls, config, cacheLoader);
    }

    /**
     * @return Stencil client for getting descriptors from classes in classpath
     */
    public static StencilClient getClient() {
        return new ClassLoadStencilClient();
    }
}
