package com.gojek.de.stencil;

import com.google.common.cache.CacheBuilder;
import com.google.common.cache.CacheLoader;
import com.google.common.cache.LoadingCache;
import com.google.protobuf.Descriptors;
import org.apache.commons.lang3.StringUtils;

import java.io.Serializable;
import java.time.Duration;
import java.util.Map;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

public class URLStencilClient extends StencilClient implements Serializable {
    private String url;
    LoadingCache<String, Map<String, Descriptors.Descriptor>> descriptorCache;
    public static final Duration DEFAULT_TTL = Duration.ofMinutes(30L);


    public Descriptors.Descriptor get(String className) {
        try {
            return descriptorCache.get(url).get(className);
        } catch (ExecutionException e) {
            throw new StencilRuntimeException(e);
        }
    }


    public URLStencilClient(String url, Map<String, String> config, CacheLoader cacheLoader) {
        Duration ttl = StringUtils.isBlank(config.get("TTL_IN_MINUTES")) ?
                DEFAULT_TTL : Duration.ofMinutes(Long.parseLong(config.get("TTL_IN_MINUTES")));
        this.url = url;
        descriptorCache = CacheBuilder.newBuilder()
                .refreshAfterWrite(10L, TimeUnit.SECONDS)
                .build(cacheLoader);
    }


}
