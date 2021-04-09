package io.odpf.stencil.config;

import lombok.Builder;
import lombok.Getter;

@Getter
@Builder
public class StencilConfig {
    @Builder.Default
    Integer fetchTimeoutMs = 10000;
    @Builder.Default
    Integer fetchRetries = 4;
    @Builder.Default
    Long fetchBackoffMinMs = 0L;
    @Builder.Default
    Boolean cacheAutoRefresh = false;
    @Builder.Default
    Long cacheTtlMs = 0L;
}
