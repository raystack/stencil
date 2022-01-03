package io.odpf.stencil.config;

import lombok.Builder;
import lombok.Getter;
import java.util.concurrent.TimeUnit;

@Getter
@Builder
public class StencilConfig {

    /**
     * HTTP timeout while fetching a protobuf descriptor set file from remote URL
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
     * Min starting Backoff when retrying, backoff gets doubled every retry till {@link #fetchRetries} are reached
     * @param fetchBackoffMinMs Min starting Backoff when retrying, backoff gets doubled every retry till {@link #fetchRetries} are reached
     * @return Min starting Backoff when retrying, backoff gets doubled every retry till {@link #fetchRetries} are reached
     */
    @Builder.Default
    Long fetchBackoffMinMs = 0L;
    /**
     * Bearer token to be used in HTTP Authorization header while fetching protobuf descriptor set files
     * @param fetchAuthBearerToken Bearer token to be used in HTTP Authorization header while fetching protobuf descriptor set files
     * @return Bearer token to be used in HTTP Authorization header while fetching protobuf descriptor set files
     */
    String fetchAuthBearerToken;
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
}
