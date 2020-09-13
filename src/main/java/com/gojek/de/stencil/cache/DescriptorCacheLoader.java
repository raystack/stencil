package com.gojek.de.stencil.cache;

import com.gojek.de.stencil.DescriptorMapBuilder;
import com.gojek.de.stencil.exception.StencilRuntimeException;
import com.gojek.de.stencil.http.RemoteFile;
import com.gojek.de.stencil.models.DescriptorAndTypeName;
import com.google.common.cache.CacheLoader;
import com.google.common.util.concurrent.Futures;
import com.google.common.util.concurrent.ListenableFuture;
import com.google.common.util.concurrent.ListenableFutureTask;
import com.google.protobuf.Descriptors;
import com.timgroup.statsd.StatsDClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayInputStream;
import java.io.Closeable;
import java.io.IOException;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class DescriptorCacheLoader extends CacheLoader<String, Map<String, DescriptorAndTypeName>> implements Closeable {
    private static final Integer DEFAULT_THREAD_POOL = 2;
    private final Logger logger = LoggerFactory.getLogger(DescriptorCacheLoader.class);
    private StatsDClient statsDClient;
    private ExecutorService executor = Executors.newFixedThreadPool(DEFAULT_THREAD_POOL);
    private RemoteFile remoteFile;
    private ProtoUpdateListener protoUpdateListener;
    private boolean asyncRefresh;

    public DescriptorCacheLoader(RemoteFile remoteFile, StatsDClient statsDClient, ProtoUpdateListener protoUpdateListener, boolean asyncRefresh) {
        this.remoteFile = remoteFile;
        this.statsDClient = statsDClient;
        this.protoUpdateListener = protoUpdateListener;
        this.asyncRefresh = asyncRefresh;
    }

    @Override
    public Map<String, DescriptorAndTypeName> load(String key) {
        logger.info("loading stencil cache");
        return refreshMap(key, new HashMap<>());
    }

    @Override
    public ListenableFuture<Map<String, DescriptorAndTypeName>> reload(final String key, final Map<String, DescriptorAndTypeName> prevDescriptor) {
        logger.info("reloading the cache to get the new descriptors");
        if(!asyncRefresh) {
            return Futures.immediateFuture(refreshMap(key, prevDescriptor));
        }
        ListenableFutureTask<Map<String, DescriptorAndTypeName>> task = ListenableFutureTask.create(
                () -> {
                    try {
                        return refreshMap(key, prevDescriptor);
                    } catch (Throwable e) {
                        logger.info("Exception on refreshing stencil descriptor", e);
                        return prevDescriptor;
                    }
                }
        );
        executor.execute(task);
        return task;
    }

    private Map<String, DescriptorAndTypeName> refreshMap(String url, final Map<String, DescriptorAndTypeName> prevDescriptor) {
        try {
            logger.info("fetching descriptors from {}", url);
            byte[] descriptorBin = remoteFile.fetch(url);
            logger.info("successfully fetched {}", url);
            InputStream inputStream = new ByteArrayInputStream(descriptorBin);
            statsDClient.count("stencil.client.refresh,status=success", 1);
            Map<String, DescriptorAndTypeName> newDescriptorsMap = new DescriptorMapBuilder().buildFrom(inputStream);

            if (protoUpdateListener != null && !prevDescriptor.isEmpty()) {
                protoUpdateListener.onProtoUpdate(url, newDescriptorsMap);
            }

            return newDescriptorsMap;
        } catch (IOException | Descriptors.DescriptorValidationException e) {
            statsDClient.count("stencil.client.refresh,status=failed", 1);
            throw new StencilRuntimeException(e);
        }
    }

    @Override
    public void close() throws IOException {
        remoteFile.close();
        executor.shutdown();
    }
}

