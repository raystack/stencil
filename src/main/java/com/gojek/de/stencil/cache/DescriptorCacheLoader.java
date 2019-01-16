package com.gojek.de.stencil.cache;

import com.gojek.de.stencil.DescriptorMapBuilder;
import com.gojek.de.stencil.RemoteFile;
import com.gojek.de.stencil.StencilRuntimeException;
import com.google.common.cache.CacheLoader;
import com.google.common.util.concurrent.ListenableFuture;
import com.google.common.util.concurrent.ListenableFutureTask;
import com.google.protobuf.Descriptors;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Map;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class DescriptorCacheLoader extends CacheLoader<String, Map<String, Descriptors.Descriptor>> {
    public static final String DEFAULT_STENCIL_TIMEOUT_MS = "10000";
    public static final String DEFAULT_STENCIL_BACKOFF_MS = "1000";
    public static final String DEFAULT_STENCIL_RETRIES = "4";
    public static final Integer DEFAULT_THREAD_POOL = 10;

    ExecutorService executor = Executors.newFixedThreadPool(DEFAULT_THREAD_POOL);
    Map<String, String> config;
    RemoteFile remoteFile;
    final Logger logger = LoggerFactory.getLogger(DescriptorCacheLoader.class);


    public DescriptorCacheLoader(Map<String, String> config, RemoteFile remoteFile) {
        this.config = config;
        this.remoteFile = remoteFile;
    }


    @Override
    public Map<String, Descriptors.Descriptor> load(String key) {
        return refreshMap(key, config);
    }

    @Override
    public ListenableFuture<Map<String, Descriptors.Descriptor>> reload(String key, Map<String, Descriptors.Descriptor> prevDescriptor) {
        ListenableFutureTask<Map<String, Descriptors.Descriptor>> task = ListenableFutureTask.create(
                () -> refreshMap(key, config)
        );
        executor.execute(task);
        return task;
    }


    private Map<String, Descriptors.Descriptor> refreshMap(String url, Map<String, String> config) {
        int timeout = Integer.parseInt(StringUtils.isBlank(config.get("STENCIL_TIMEOUT_MS")) ?
                DEFAULT_STENCIL_TIMEOUT_MS : config.get("STENCIL_TIMEOUT_MS"));
        int backoffMs = Integer.parseInt(StringUtils.isBlank(config.get("STENCIL_BACKOFF_MS")) ?
                DEFAULT_STENCIL_BACKOFF_MS : config.get("STENCIL_BACKOFF_MS"));
        int retries = Integer.parseInt(StringUtils.isBlank(config.get("STENCIL_RETRIES")) ?
                DEFAULT_STENCIL_RETRIES : config.get("STENCIL_RETRIES"));
        int retryCount = retries;


        try {
            logger.info("fetching descriptors from {} with timeout: {}ms, backoff: {}ms {} retries pending", url, timeout, backoffMs, retryCount);
            byte[] descriptorBin = remoteFile.fetch(url, timeout);
            logger.info("successfully fetched {}", url);
            InputStream inputStream = new ByteArrayInputStream(descriptorBin);
            return new DescriptorMapBuilder().buildFrom(inputStream);

        } catch (IOException | Descriptors.DescriptorValidationException e) {
            throw new StencilRuntimeException(e);
        }
    }
}

