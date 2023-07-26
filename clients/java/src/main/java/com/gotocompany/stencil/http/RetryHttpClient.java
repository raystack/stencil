package com.gotocompany.stencil.http;

import com.gotocompany.stencil.config.StencilConfig;
import org.apache.hc.client5.http.HttpRequestRetryStrategy;
import org.apache.hc.client5.http.config.ConnectionConfig;
import org.apache.hc.client5.http.config.RequestConfig;
import org.apache.hc.client5.http.cookie.StandardCookieSpec;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.client5.http.impl.classic.HttpClientBuilder;
import org.apache.hc.client5.http.impl.io.PoolingHttpClientConnectionManager;
import org.apache.hc.client5.http.impl.io.PoolingHttpClientConnectionManagerBuilder;
import org.apache.hc.core5.http.HttpRequest;
import org.apache.hc.core5.http.HttpResponse;
import org.apache.hc.core5.http.io.SocketConfig;
import org.apache.hc.core5.http.protocol.HttpContext;
import org.apache.hc.core5.pool.PoolConcurrencyPolicy;
import org.apache.hc.core5.pool.PoolReusePolicy;
import org.apache.hc.core5.util.TimeValue;
import org.apache.hc.core5.util.Timeout;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.concurrent.TimeUnit;

public class RetryHttpClient {
    private static final Logger logger = LoggerFactory.getLogger(RemoteFileImpl.class);


    public static CloseableHttpClient create(StencilConfig stencilConfig) {

        int timeout = stencilConfig.getFetchTimeoutMs();
        long backoffMs = stencilConfig.getFetchBackoffMinMs();
        int retries = stencilConfig.getFetchRetries();
        Timeout defaultTimeout = Timeout.of(timeout, TimeUnit.MILLISECONDS);
        PoolingHttpClientConnectionManager connectionManager = PoolingHttpClientConnectionManagerBuilder.create()
                .setDefaultSocketConfig(SocketConfig.custom()
                        .setSoTimeout(defaultTimeout)
                        .build())
                .setPoolConcurrencyPolicy(PoolConcurrencyPolicy.STRICT)
                .setConnPoolPolicy(PoolReusePolicy.LIFO)
                .setDefaultConnectionConfig(ConnectionConfig.custom()
                        .setSocketTimeout(defaultTimeout)
                        .setConnectTimeout(defaultTimeout)
                        .setTimeToLive(TimeValue.ofMinutes(10))
                        .build())
                .build();
        logger.info("initialising HTTP client with timeout: {}ms, backoff: {}ms, max retry attempts: {}", timeout, backoffMs, retries);

        return HttpClientBuilder.create()
                .setDefaultRequestConfig(RequestConfig.custom()
                        .setCookieSpec(StandardCookieSpec.STRICT)
                        .build())
                .setConnectionManager(connectionManager)
                .setDefaultHeaders(stencilConfig.getFetchHeaders())
                .setConnectionManagerShared(true)
                .setRetryStrategy(new HttpRequestRetryStrategy() {
                    private TimeValue waitPeriod = TimeValue.of(backoffMs, TimeUnit.MILLISECONDS);

                    @Override
                    public boolean retryRequest(HttpRequest request, IOException exception, int execCount, HttpContext context) {
                        return false;
                    }

                    @Override
                    public boolean retryRequest(HttpResponse response, int execCount, HttpContext context) {
                        if (execCount <= retries && response.getCode() >= 400) {
                            logger.info("Retrying requests, attempts left: {}", retries - execCount);
                            waitPeriod = TimeValue.of(waitPeriod.toMilliseconds() * 2, TimeUnit.MILLISECONDS);
                            return true;
                        }
                        return false;
                    }

                    @Override
                    public TimeValue getRetryInterval(HttpResponse response, int execCount, HttpContext context) {
                        return waitPeriod;
                    }
                })
                .build();
    }

}
