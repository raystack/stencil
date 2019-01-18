package com.gojek.de.stencil.http;

import org.apache.commons.lang3.StringUtils;
import org.apache.http.HttpResponse;
import org.apache.http.client.ServiceUnavailableRetryStrategy;
import org.apache.http.client.config.RequestConfig;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.protocol.HttpContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;

public class RetryHttpClient {
    private final Logger logger = LoggerFactory.getLogger(RemoteFileImpl.class);


    private static final String DEFAULT_STENCIL_TIMEOUT_MS = "10000";
    private static final String DEFAULT_STENCIL_BACKOFF_MS = "2000";
    private static final String DEFAULT_STENCIL_RETRIES = "4";

    public CloseableHttpClient create(Map<String, String> config) {
        int timeout = Integer.parseInt(StringUtils.isBlank(config.get("STENCIL_TIMEOUT_MS")) ?
                DEFAULT_STENCIL_TIMEOUT_MS : config.get("STENCIL_TIMEOUT_MS"));
        int backoffMs = Integer.parseInt(StringUtils.isBlank(config.get("STENCIL_BACKOFF_MS")) ?
                DEFAULT_STENCIL_BACKOFF_MS : config.get("STENCIL_BACKOFF_MS"));
        int retries = Integer.parseInt(StringUtils.isBlank(config.get("STENCIL_RETRIES")) ?
                DEFAULT_STENCIL_RETRIES : config.get("STENCIL_RETRIES"));

        logger.info("initialising HTTP client with timeout: {}ms, backoff: {}ms, max retry attempts: {}", timeout, backoffMs, retries);


        RequestConfig requestConfig = RequestConfig.custom()
                .setConnectTimeout(timeout)
                .setSocketTimeout(timeout).build();


        return HttpClientBuilder.create()
                .setDefaultRequestConfig(requestConfig)
                .setConnectionManagerShared(true)
                .setServiceUnavailableRetryStrategy(new ServiceUnavailableRetryStrategy() {
                    int waitPeriod = backoffMs;

                    @Override
                    public boolean retryRequest(HttpResponse response, int executionCount, HttpContext context) {
                        if (executionCount <= retries && response.getStatusLine().getStatusCode() >= 400) {
                            logger.info("Retrying requests, attempts left: {}", retries - executionCount);
                            waitPeriod *= 2;
                            return true;
                        }
                        return false;
                    }

                    @Override
                    public long getRetryInterval() {
                        return waitPeriod;
                    }
                })
                .build();
    }

}
