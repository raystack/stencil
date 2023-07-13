package org.raystack.stencil;

import org.raystack.stencil.cache.SchemaCacheLoader;
import org.raystack.stencil.client.ClassLoadStencilClient;
import org.raystack.stencil.client.MultiURLStencilClient;
import org.raystack.stencil.client.StencilClient;
import org.raystack.stencil.client.URLStencilClient;
import org.raystack.stencil.config.StencilConfig;
import org.raystack.stencil.http.RemoteFileImpl;
import org.raystack.stencil.http.RetryHttpClient;

import java.util.List;


/**
 * Provides static methods for the creation of {@link org.raystack.stencil.client.StencilClient}
 * object with configurations and various options like
 * single URLs, multiple URLs, statsd client for monitoring
 * and {@link org.raystack.stencil.SchemaUpdateListener} for callback on schema change.
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
