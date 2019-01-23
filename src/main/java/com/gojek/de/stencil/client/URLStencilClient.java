package com.gojek.de.stencil.client;

import com.gojek.de.stencil.cache.DescriptorCacheLoader;
import com.gojek.de.stencil.exception.StencilRuntimeException;
import com.gojek.de.stencil.utils.RandomUtils;
import com.google.common.base.Ticker;
import com.google.common.cache.CacheBuilder;
import com.google.common.cache.LoadingCache;
import com.google.common.util.concurrent.UncheckedExecutionException;
import com.google.protobuf.Descriptors;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.Serializable;
import java.time.Duration;
import java.util.Map;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

import static com.google.common.base.Ticker.systemTicker;

public class URLStencilClient implements Serializable, StencilClient {
    private String url;
    private DescriptorCacheLoader cacheLoader;
    private LoadingCache<String, Map<String, Descriptors.Descriptor>> descriptorCache;
    private Duration ttl;
    private static final int DEFAULT_TTL_MIN = 30;
    private static final int DEFAULT_TTL_MAX = 60;
    private final Logger logger = LoggerFactory.getLogger(URLStencilClient.class);

    public Descriptors.Descriptor get(String className) {
        try {
            return descriptorCache.get(url).get(className);
        } catch (UncheckedExecutionException | ExecutionException e) {
            throw new StencilRuntimeException(e);
        }
    }


    public URLStencilClient(String url, Map<String, String> config, DescriptorCacheLoader cacheLoader) {
        this(url, config, cacheLoader, systemTicker());
    }


    public URLStencilClient(String url, Map<String, String> config, DescriptorCacheLoader cacheLoader, Ticker ticker) {
        this.ttl = StringUtils.isBlank(config.get("TTL_IN_MINUTES")) ?
                getDefaultTTL() : Duration.ofMinutes(Long.parseLong(config.get("TTL_IN_MINUTES")));
        this.url = url;
        this.cacheLoader = cacheLoader;
        descriptorCache = CacheBuilder.newBuilder().ticker(ticker)
                .refreshAfterWrite(ttl.toMinutes(), TimeUnit.MINUTES)
                .build(cacheLoader);
        logger.info("initialising URL Stencil client with TTL: {} minutes", ttl.toMinutes());

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
