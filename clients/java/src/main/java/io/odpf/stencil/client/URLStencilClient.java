package io.odpf.stencil.client;

import com.google.common.base.Ticker;
import com.google.common.cache.CacheBuilder;
import com.google.common.cache.LoadingCache;
import com.google.common.util.concurrent.UncheckedExecutionException;
import com.google.protobuf.Descriptors;
import io.odpf.stencil.cache.DescriptorCacheLoader;
import io.odpf.stencil.config.StencilConfig;
import io.odpf.stencil.exception.StencilRuntimeException;
import io.odpf.stencil.models.DescriptorAndTypeName;
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
    private DescriptorCacheLoader cacheLoader;
    private LoadingCache<String, Map<String, DescriptorAndTypeName>> descriptorCache;
    private long ttlMs;
    private final Logger logger = LoggerFactory.getLogger(URLStencilClient.class);
    private boolean shouldAutoRefreshCache;

    /**
     * @param url List of URLs to fetch protobuf descriptor sets from
     * @param config Stencil configs
     * @param cacheLoader Extension of Guava {@link com.google.common.cache.CacheLoader} for Proto Descriptor sets
     */
    public URLStencilClient(String url, StencilConfig config, DescriptorCacheLoader cacheLoader) {
        this(url, config, cacheLoader, systemTicker());
    }

    /**
     * @param url List of URLs to fetch protobuf descriptor sets from
     * @param stencilConfig Stencil configs
     * @param cacheLoader Extension of Guava {@link com.google.common.cache.CacheLoader} for Proto Descriptor sets
     * @param ticker Ticker to be used as time source in Guava cache
     */
    public URLStencilClient(String url, StencilConfig stencilConfig, DescriptorCacheLoader cacheLoader, Ticker ticker) {
        this.shouldAutoRefreshCache = stencilConfig.getCacheAutoRefresh();
        this.ttlMs = stencilConfig.getCacheTtlMs();
        this.url = url;
        this.cacheLoader = cacheLoader;

        if (stencilConfig.getCacheAutoRefresh()) {
            descriptorCache = CacheBuilder.newBuilder().ticker(ticker)
                    .refreshAfterWrite(ttlMs, TimeUnit.MILLISECONDS)
                    .build(cacheLoader);
            logger.info("configuring URL Stencil client with TTL: {} milliseconds", ttlMs);
        } else {
            descriptorCache = CacheBuilder.newBuilder().ticker(ticker)
                    .build(cacheLoader);
        }

        logger.info("initialising URL Stencil client with auto refresh: {}", shouldAutoRefreshCache);
    }

    /**
     * @param className Class name of the required protobuf schema
     * @return {@link com.google.protobuf.Descriptors.Descriptor} describing the schema of given class name
     */
    @Override
    public Descriptors.Descriptor get(String className) {
        try {
            DescriptorAndTypeName descriptorAndTypeName = descriptorCache.get(url).get(className);
            return descriptorAndTypeName != null ? descriptorAndTypeName.getDescriptor() : null;
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
                DescriptorAndTypeName descriptorAndTypeName = mapEntry.getValue();
                if (descriptorAndTypeName != null) {
                    descriptorMap.put(mapEntry.getKey(), descriptorAndTypeName.getDescriptor());
                }
            });
            return descriptorMap;
        } catch (UncheckedExecutionException | ExecutionException e) {
            throw new StencilRuntimeException(e);
        }
    }

    @Override
    public Map<String, String> getTypeNameToPackageNameMap() {
        try {
            Map<String, String> typeNameMap = new HashMap<>();
            descriptorCache.get(url).entrySet().stream().forEach(mapEntry -> {
                DescriptorAndTypeName descriptorAndTypeName = mapEntry.getValue();
                if (descriptorAndTypeName != null) {
                    typeNameMap.put(descriptorAndTypeName.getTypeName(), mapEntry.getKey());
                }
            });
            return typeNameMap;
        } catch (UncheckedExecutionException | ExecutionException e) {
            throw new StencilRuntimeException(e);
        }
    }

    @Override
    public Map<String, DescriptorAndTypeName> getAllDescriptorAndTypeName() {
        try {
            return descriptorCache.get(url);
        } catch (ExecutionException e) {
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

    @Override
    public boolean shouldAutoRefreshCache() {
        return shouldAutoRefreshCache;
    }
}
