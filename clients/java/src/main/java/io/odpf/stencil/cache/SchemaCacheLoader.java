package com.gotocompany.stencil.cache;

import com.google.common.cache.CacheLoader;
import com.google.common.util.concurrent.Futures;
import com.google.common.util.concurrent.ListenableFuture;
import com.google.common.util.concurrent.ListenableFutureTask;
import com.google.protobuf.Descriptors;
import com.gotocompany.stencil.config.StencilConfig;
import com.gotocompany.stencil.http.RemoteFile;
import com.timgroup.statsd.StatsDClient;
import com.gotocompany.stencil.SchemaUpdateListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.Closeable;
import java.io.IOException;
import java.util.Map;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class SchemaCacheLoader extends CacheLoader<String, Map<String, Descriptors.Descriptor>> implements Closeable {
    private static final Integer DEFAULT_THREAD_POOL = 2;
    private final Logger logger = LoggerFactory.getLogger(SchemaCacheLoader.class);
    private StatsDClient statsDClient;
    private ExecutorService executor = Executors.newFixedThreadPool(DEFAULT_THREAD_POOL);
    private RemoteFile remoteFile;
    private SchemaUpdateListener protoUpdateListener;
    private boolean shouldRefresh;
    private SchemaRefreshStrategy refreshStrategy;

    public SchemaCacheLoader(RemoteFile remoteFile, StencilConfig config) {
        this.remoteFile = remoteFile;
        this.statsDClient = config.getStatsDClient();
        this.protoUpdateListener = config.getUpdateListener();
        this.shouldRefresh = config.getCacheAutoRefresh();
        this.refreshStrategy = config.getRefreshStrategy();
    }

    @Override
    public Map<String, Descriptors.Descriptor> load(String key) {
        logger.info("loading stencil cache");
        return refreshStrategy.refresh(key, remoteFile, null);
    }

    @Override
    public ListenableFuture<Map<String, Descriptors.Descriptor>> reload(final String key, final Map<String, Descriptors.Descriptor> prevDescriptor) {
        if (!shouldRefresh) {
            return Futures.immediateFuture(prevDescriptor);
        }
        logger.info("reloading the cache to get the new descriptors");
        ListenableFutureTask<Map<String, Descriptors.Descriptor>> task = ListenableFutureTask.create(() -> {
            try {
                Map<String, Descriptors.Descriptor> newDescriptor = refreshStrategy.refresh(key, remoteFile, prevDescriptor);
                statsDClient.count("stencil.client.refresh,status=success", 1);
                if (prevDescriptor != newDescriptor && protoUpdateListener != null) {
                    protoUpdateListener.onSchemaUpdate(newDescriptor);
                }
                return newDescriptor;
            } catch (Throwable e) {
                statsDClient.count("stencil.client.refresh,status=failed", 1);
                logger.info("Exception on refreshing stencil descriptor", e);
                return prevDescriptor;
            }
        });
        executor.execute(task);
        return task;
    }

    @Override
    public void close() throws IOException {
        remoteFile.close();
        executor.shutdown();
    }
}
