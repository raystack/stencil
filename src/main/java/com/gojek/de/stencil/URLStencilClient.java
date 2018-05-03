package com.gojek.de.stencil;

import com.google.protobuf.Descriptors;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;

import java.io.IOException;
import java.io.InputStream;
import java.util.Map;

@Slf4j
public class URLStencilClient implements StencilClient {

    public static final String DEFAULT_STENCIL_TIMEOUT_MS = "10000";
    public static final String DEFAULT_STENCIL_BACKOFF_MS = "1000";
    public static final String DEFAULT_STENCIL_RETRIES = "4";

    public URLStencilClient(String url, Map<String, String> options) {
        this.options = options;
        this.url = url;
    }

    private String url;
    private Map<String, String> options;
    private Map<String, Descriptors.Descriptor> descriptorMap;

    public Descriptors.Descriptor get(String className) {
        if (descriptorMap == null) {
            load();
        }
        return descriptorMap.get(className);
    }

    @Override
    public void load() {
        int timeout = Integer.parseInt(StringUtils.isBlank(options.get("STENCIL_TIMEOUT_MS")) ?
                DEFAULT_STENCIL_TIMEOUT_MS : options.get("STENCIL_TIMEOUT_MS"));
        int backoffMs = Integer.parseInt(StringUtils.isBlank(options.get("STENCIL_BACKOFF_MS")) ?
                DEFAULT_STENCIL_BACKOFF_MS : options.get("STENCIL_BACKOFF_MS"));
        int retries = Integer.parseInt(StringUtils.isBlank(options.get("STENCIL_RETRIES")) ?
                DEFAULT_STENCIL_RETRIES : options.get("STENCIL_RETRIES"));
        int retryCount = retries;
        //get schema from server
        while (true) {
            try {
                log.info("fetching descriptors from {} with timeout: {}ms, backoff: {}ms {} retries pending", url);
                InputStream is = new RemoteFile().fetch(url, timeout);
                descriptorMap = new DescriptorMapBuilder().buildFrom(is);
                break;
            } catch (IOException | RuntimeException | Descriptors.DescriptorValidationException e) {
                if (retryCount < 1) {
                    throw new StencilRuntimeException(e);
                }
                log.error(e.getMessage());
            }
            retryCount--;
            try {
                Thread.sleep(backoffMs * (retries - retryCount));
            } catch (InterruptedException e) {
                throw new StencilRuntimeException(e);
            }
        }

    }

    @Override
    public void reload() {
        load();
    }
}
