package io.odpf.stencil.http;

import io.odpf.stencil.config.StencilConfig;
import org.apache.http.HttpResponse;
import org.apache.http.client.ServiceUnavailableRetryStrategy;
import org.apache.http.client.config.RequestConfig;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.protocol.HttpContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class RetryHttpClient {
    private final Logger logger = LoggerFactory.getLogger(RemoteFileImpl.class);


    public CloseableHttpClient create(StencilConfig stencilConfig) {

        int timeout = stencilConfig.getFetchTimeoutMs();
        long backoffMs = stencilConfig.getFetchBackoffMinMs();
        int retries = stencilConfig.getFetchRetries();

        logger.info("initialising HTTP client with timeout: {}ms, backoff: {}ms, max retry attempts: {}", timeout, backoffMs, retries);

        RequestConfig requestConfig = RequestConfig.custom()
                .setConnectTimeout(timeout)
                .setSocketTimeout(timeout).build();

        return HttpClientBuilder.create()
                .setDefaultRequestConfig(requestConfig)
                .setDefaultHeaders(stencilConfig.getFetchHeaders())
                .setConnectionManagerShared(true)
                .setServiceUnavailableRetryStrategy(new ServiceUnavailableRetryStrategy() {
                    long waitPeriod = backoffMs;

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
