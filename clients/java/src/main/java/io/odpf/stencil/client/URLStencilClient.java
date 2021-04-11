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
import io.odpf.stencil.utils.RandomUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.Serializable;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

import static com.google.common.base.Ticker.systemTicker;

public class URLStencilClient implements Serializable, StencilClient {
    private String url;
    private DescriptorCacheLoader cacheLoader;
    private LoadingCache<String, Map<String, DescriptorAndTypeName>> descriptorCache;
    private Duration ttl;
    private static final long DEFAULT_TTL_MIN = TimeUnit.MINUTES.toMillis(30);
    private static final long DEFAULT_TTL_MAX = TimeUnit.MINUTES.toMillis(60);
    private final Logger logger = LoggerFactory.getLogger(URLStencilClient.class);
    private boolean shouldAutoRefreshCache;

    public URLStencilClient(String url, StencilConfig config, DescriptorCacheLoader cacheLoader) {
        this(url, config, cacheLoader, systemTicker());
    }

    public URLStencilClient(String url, StencilConfig stencilConfig, DescriptorCacheLoader cacheLoader, Ticker ticker) {
        this.shouldAutoRefreshCache = stencilConfig.getCacheAutoRefresh();
        this.ttl = stencilConfig.getCacheTtlMs() != 0 ? Duration.ofMillis(stencilConfig.getCacheTtlMs()) :
                getDefaultTTL();
        this.url = url;
        this.cacheLoader = cacheLoader;

        if (stencilConfig.getCacheAutoRefresh()) {
            descriptorCache = CacheBuilder.newBuilder().ticker(ticker)
                    .refreshAfterWrite(ttl.toMillis(), TimeUnit.MILLISECONDS)
                    .build(cacheLoader);
            logger.info("configuring URL Stencil client with TTL: {} milliseconds", ttl.toMillis());
        } else {
            descriptorCache = CacheBuilder.newBuilder().ticker(ticker)
                    .build(cacheLoader);
        }

        logger.info("initialising URL Stencil client with auto refresh: {}", shouldAutoRefreshCache);
    }

    @Override
    public Descriptors.Descriptor get(String className) {
        try {
            DescriptorAndTypeName descriptorAndTypeName = descriptorCache.get(url).get(className);
            return descriptorAndTypeName != null ? descriptorAndTypeName.getDescriptor() : null;
        } catch (UncheckedExecutionException | ExecutionException e) {
            throw new StencilRuntimeException(e);
        }
    }

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

    public Duration getTTL() {
        return ttl;
    }

    private Duration getDefaultTTL() {
        long randomNumberInRange = new RandomUtils().getRandomNumberInRange(DEFAULT_TTL_MIN, DEFAULT_TTL_MAX);
        return Duration.ofMillis(randomNumberInRange);

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
