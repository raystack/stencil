package io.odpf.stencil.client;

import com.google.common.base.Ticker;
import com.google.common.cache.CacheBuilder;
import com.google.common.cache.LoadingCache;
import com.google.common.util.concurrent.UncheckedExecutionException;
import com.google.protobuf.Descriptors;
import io.odpf.stencil.cache.SchemaCacheLoader;
import io.odpf.stencil.config.StencilConfig;
import io.odpf.stencil.exception.StencilRuntimeException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.Serializable;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

import static com.google.common.base.Ticker.systemTicker;

/**
 * {@link StencilClient} implementation that can fetch descriptor sets from single URL
 */
public class URLStencilClient implements Serializable, StencilClient {
    private String url;
    private SchemaCacheLoader cacheLoader;
    private LoadingCache<String, Map<String, Descriptors.Descriptor>> descriptorCache;
    private long ttlMs;
    private final Logger logger = LoggerFactory.getLogger(URLStencilClient.class);

    /**
     * @param url List of URLs to fetch protobuf descriptor sets from
     * @param config Stencil configs
     * @param cacheLoader Extension of Guava {@link com.google.common.cache.CacheLoader} for Proto Descriptor sets
     */
    public URLStencilClient(String url, StencilConfig config, SchemaCacheLoader cacheLoader) {
        this(url, config, cacheLoader, systemTicker());
    }

    /**
     * @param url List of URLs to fetch protobuf descriptor sets from
     * @param stencilConfig Stencil configs
     * @param cacheLoader Extension of Guava {@link com.google.common.cache.CacheLoader} for Proto Descriptor sets
     * @param ticker Ticker to be used as time source in Guava cache
     */
    public URLStencilClient(String url, StencilConfig stencilConfig, SchemaCacheLoader cacheLoader, Ticker ticker) {
        this.ttlMs = stencilConfig.getCacheTtlMs();
        this.url = url;
        this.cacheLoader = cacheLoader;

        descriptorCache = CacheBuilder.newBuilder().ticker(ticker).refreshAfterWrite(ttlMs, TimeUnit.MILLISECONDS).build(cacheLoader);
        logger.info("configuring URL Stencil client with TTL: {} milliseconds, auto refresh: {}", ttlMs, stencilConfig.getCacheAutoRefresh());
    }

    /**
     * @param className Class name of the required protobuf schema
     * @return {@link com.google.protobuf.Descriptors.Descriptor} describing the schema of given class name
     */
    @Override
    public Descriptors.Descriptor get(String className) {
        try {
            return descriptorCache.get(url).get(className);
        } catch (UncheckedExecutionException | ExecutionException e) {
            throw new StencilRuntimeException(e);
        }
    }

    /**
     * @return Get a map containing all loaded protobuf schema names and their descriptors
     */
    @Override
    public Map<String, Descriptors.Descriptor> getAll() {
        try {
            Map<String, Descriptors.Descriptor> descriptorMap = new HashMap<>();
            descriptorCache.get(url).entrySet().stream().forEach(mapEntry -> {
                Descriptors.Descriptor desc = mapEntry.getValue();
                if (desc != null) {
                    descriptorMap.put(mapEntry.getKey(), desc);
                }
            });
            return descriptorMap;
        } catch (UncheckedExecutionException | ExecutionException e) {
            throw new StencilRuntimeException(e);
        }
    }

    public void refresh() {
        descriptorCache.refresh(url);
    }

    public long getTTLMs() {
        return ttlMs;
    }

    @Override
    public void close() throws IOException {
        cacheLoader.close();
    }
}
