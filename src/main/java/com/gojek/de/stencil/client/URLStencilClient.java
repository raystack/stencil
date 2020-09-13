package com.gojek.de.stencil.client;

import com.gojek.de.stencil.cache.DescriptorCacheLoader;
import com.gojek.de.stencil.config.StencilConfig;
import com.gojek.de.stencil.exception.StencilRuntimeException;
import com.gojek.de.stencil.models.DescriptorAndTypeName;
import com.gojek.de.stencil.utils.RandomUtils;
import com.google.common.base.Ticker;
import com.google.common.cache.CacheBuilder;
import com.google.common.cache.LoadingCache;
import com.google.common.util.concurrent.UncheckedExecutionException;
import com.google.protobuf.Descriptors;
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
    private static final int DEFAULT_TTL_MIN = 30;
    private static final int DEFAULT_TTL_MAX = 60;
    private final Logger logger = LoggerFactory.getLogger(URLStencilClient.class);

    public URLStencilClient(String url, StencilConfig config, DescriptorCacheLoader cacheLoader) {
        this(url, config, cacheLoader, systemTicker());
    }

    public URLStencilClient(String url, StencilConfig stencilConfig, DescriptorCacheLoader cacheLoader, Ticker ticker) {

        this.ttl = stencilConfig.getTilInMinutes() != 0 ? Duration.ofMinutes(stencilConfig.getTilInMinutes()) :
                getDefaultTTL();
        this.url = url;
        this.cacheLoader = cacheLoader;

        if (stencilConfig.shouldAutoRefreshCache()) {
            descriptorCache = CacheBuilder.newBuilder().ticker(ticker)
                    .refreshAfterWrite(ttl.toMinutes(), TimeUnit.MINUTES)
                    .build(cacheLoader);
        } else {
            descriptorCache = CacheBuilder.newBuilder().ticker(ticker)
                    .build(cacheLoader);
        }

        logger.info("initialising URL Stencil client with TTL: {} minutes", ttl.toMinutes());
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
        int randomNumberInRange = new RandomUtils().getRandomNumberInRange(DEFAULT_TTL_MIN, DEFAULT_TTL_MAX);
        return Duration.ofMinutes(randomNumberInRange);

    }

    @Override
    public void close() throws IOException {
        cacheLoader.close();
    }
}
