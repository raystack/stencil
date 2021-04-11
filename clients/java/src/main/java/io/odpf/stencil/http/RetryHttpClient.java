package io.odpf.stencil.http;

import io.odpf.stencil.config.StencilConfig;
import io.odpf.stencil.utils.RandomUtils;
import org.apache.http.Header;
import org.apache.http.HttpHeaders;
import org.apache.http.HttpResponse;
import org.apache.http.client.ServiceUnavailableRetryStrategy;
import org.apache.http.client.config.RequestConfig;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.message.BasicHeader;
import org.apache.http.protocol.HttpContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;

public class RetryHttpClient {
    private final Logger logger = LoggerFactory.getLogger(RemoteFileImpl.class);

    private static final int DEFAULT_STENCIL_BACKOFF_MS_MIN = 2000;
    private static final int DEFAULT_STENCIL_BACKOFF_MS_MAX = 5000;

    public CloseableHttpClient create(StencilConfig stencilConfig) {

        int timeout = stencilConfig.getFetchTimeoutMs();
        long backoffMs = stencilConfig.getFetchBackoffMinMs() != 0 ? stencilConfig.getFetchBackoffMinMs() :
                new RandomUtils().getRandomNumberInRange(DEFAULT_STENCIL_BACKOFF_MS_MIN, DEFAULT_STENCIL_BACKOFF_MS_MAX);
        int retries = stencilConfig.getFetchRetries();
        List<Header> defaultHeaders = new ArrayList<>();

        if (stencilConfig.getFetchAuthBearerToken() != null) {
            String authHeaderValue = "Bearer " + stencilConfig.getFetchAuthBearerToken();
            defaultHeaders.add(new BasicHeader(HttpHeaders.AUTHORIZATION, authHeaderValue));
        }

        logger.info("initialising HTTP client with timeout: {}ms, backoff: {}ms, max retry attempts: {}", timeout, backoffMs, retries);


        RequestConfig requestConfig = RequestConfig.custom()
                .setConnectTimeout(timeout)
                .setSocketTimeout(timeout).build();


        return HttpClientBuilder.create()
                .setDefaultRequestConfig(requestConfig)
                .setDefaultHeaders(defaultHeaders)
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
