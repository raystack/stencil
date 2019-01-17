package com.gojek.de.stencil.client;

import com.gojek.de.stencil.exception.StencilRuntimeException;
import com.google.common.base.Ticker;
import com.google.common.cache.CacheBuilder;
import com.google.common.cache.CacheLoader;
import com.google.common.cache.LoadingCache;
import com.google.common.util.concurrent.UncheckedExecutionException;
import com.google.protobuf.Descriptors;
import org.apache.commons.lang3.StringUtils;

import java.io.Serializable;
import java.time.Duration;
import java.util.Map;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

import static com.google.common.base.Ticker.systemTicker;

public class URLStencilClient implements Serializable, StencilClient {
    private String url;
    private LoadingCache<String, Map<String, Descriptors.Descriptor>> descriptorCache;
    private static final Duration DEFAULT_TTL = Duration.ofMinutes(30L);


    public Descriptors.Descriptor get(String className) {
        try {
            return descriptorCache.get(url).get(className);
        } catch (UncheckedExecutionException | ExecutionException e) {
            throw new StencilRuntimeException(e);
        }
    }


    public URLStencilClient(String url, Map<String, String> config, CacheLoader cacheLoader) {
        this(url, config, cacheLoader, systemTicker());
    }


    public URLStencilClient(String url, Map<String, String> config, CacheLoader cacheLoader, Ticker ticker) {
        Duration ttl = StringUtils.isBlank(config.get("TTL_IN_MINUTES")) ?
                DEFAULT_TTL : Duration.ofMinutes(Long.parseLong(config.get("TTL_IN_MINUTES")));
        this.url = url;
        descriptorCache = CacheBuilder.newBuilder().ticker(ticker)
                .refreshAfterWrite(ttl.toMinutes(), TimeUnit.MINUTES)
                .build(cacheLoader);
    }


}
