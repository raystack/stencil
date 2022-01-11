package io.odpf.stencil.config;

import lombok.Builder;
import lombok.Getter;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;
import com.timgroup.statsd.StatsDClient;
import com.timgroup.statsd.NoOpStatsDClient;
import org.apache.http.Header;
import io.odpf.stencil.SchemaUpdateListener;
import io.odpf.stencil.cache.SchemaRefreshStrategy;

@Getter
@Builder(toBuilder = true)
public class StencilConfig {

    /**
     * HTTP timeout while fetching a protobuf descriptor set file from remote URL. Default 10000 ms
     * @param fetchTimeoutMs HTTP timeout while fetching a protobuf descriptor set file from remote URL
     * @return HTTP timeout while fetching a protobuf descriptor set file from remote URL
     */
    @Builder.Default
    Integer fetchTimeoutMs = 10000;
    /**
     * HTTP retries while fetching a protobuf descriptor set file from remote URL where returned status code is greater than equal to 400
     * @param fetchRetries HTTP retries while fetching a protobuf descriptor set file from remote URL where returned status code is greater than equal to 400
     * @return HTTP retries while fetching a protobuf descriptor set file from remote URL where returned status code is greater than equal to 400
     */
    @Builder.Default
    Integer fetchRetries = 4;
    /**
     * Min starting Backoff when retrying, backoff gets doubled every retry till {@link #fetchRetries} are reached. Default value set to 5000
     * @param fetchBackoffMinMs Min starting Backoff when retrying, backoff gets doubled every retry till {@link #fetchRetries} are reached
     * @return Min starting Backoff when retrying, backoff gets doubled every retry till {@link #fetchRetries} are reached
     */
    @Builder.Default
    Long fetchBackoffMinMs = 5000L;

    /**
     * These headers will be added to fetch request. Default empty list
     * @param fetchHeaders list of headers passed to fetch request
     * @return List of headers
     */
    @Builder.Default
    List<Header> fetchHeaders = new ArrayList<Header>();

    /**
     * enable or disable cache auto refresh, enabling would refresh the cache after {@link #cacheTtlMs} is reached
     * @param cacheAutoRefresh enable or disable cache auto refresh, enabling would refresh the cache after {@link #cacheTtlMs} is reached
     * @return will cache refresh after {@link #cacheTtlMs} is reached
     */
    @Builder.Default
    Boolean cacheAutoRefresh = false;
    /**
     * Time to Live for cache storing protobuf descriptors in milliseconds. Default duration is 24 hours
     * @param cacheTtlMs Time to Live for cache storing protobuf descriptors in milliseconds
     * @return Time to Live for cache storing protobuf descriptors in milliseconds
     */
    @Builder.Default
    Long cacheTtlMs = TimeUnit.HOURS.toMillis(24);

    /**
     * Strategy implementation of when schema needs to be refreshed. Default {@link io.odpf.stencil.cache.SchemaRefreshStrategy#longPollingStrategy}
     * @param refreshStrategy schema refresh strategy implementation
     * @return schema refresh strategy
     */
    @Builder.Default
    SchemaRefreshStrategy refreshStrategy = SchemaRefreshStrategy.longPollingStrategy();

    /**
     * updateListener will be called on new schema load. Default null
     * @param updateListener This is callback method on new schema load
     * @return schema update listener
     */
    @Builder.Default
    SchemaUpdateListener updateListener = null;

    /**
     * statsD client to capture metrics provided by stencil. Default {@link com.timgroup.statsd.NoOpStatsDClient}
     * @param statsDClient statsDClient
     * @return instance of statsDClient
     */
    @Builder.Default
    StatsDClient statsDClient = new NoOpStatsDClient();
}
