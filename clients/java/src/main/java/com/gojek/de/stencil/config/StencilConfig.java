package com.gojek.de.stencil.config;

import org.aeonbits.owner.Config;

public interface StencilConfig extends Config{
    @Key("STENCIL_TIMEOUT_MS")
    @DefaultValue("10000")
    Integer getStencilTimeoutMs();

    @Key("STENCIL_RETRIES")
    @DefaultValue("4")
    Integer getStencilRetries();

    @Key("STENCIL_BACKOFF_MS_MIN")
    @DefaultValue("0")
    Integer getStencilBackoff();

    @Key("REFRESH_CACHE")
    @DefaultValue("false")
    Boolean shouldAutoRefreshCache();

    @Key("TIL_IN_MINUTES")
    @DefaultValue("0")
    Long getTilInMinutes();
}
