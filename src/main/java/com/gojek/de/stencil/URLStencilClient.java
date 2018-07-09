package com.gojek.de.stencil;

import com.google.protobuf.Descriptors;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.Serializable;
import java.util.Map;

public class URLStencilClient extends StencilClient implements Serializable{

    final Logger log = LoggerFactory.getLogger(URLStencilClient.class);

    public static final String DEFAULT_STENCIL_TIMEOUT_MS = "10000";
    public static final String DEFAULT_STENCIL_BACKOFF_MS = "1000";
    public static final String DEFAULT_STENCIL_RETRIES = "4";

    private byte[] descriptorBin;
    private transient Map<String, Descriptors.Descriptor> descriptorMap;

    public Descriptors.Descriptor get(String className) {
        if (descriptorMap == null) {
            generateMap();
        }
        return descriptorMap.get(className);
    }


    public URLStencilClient(String url, Map<String, String> config) {
        int timeout = Integer.parseInt(StringUtils.isBlank(config.get("STENCIL_TIMEOUT_MS")) ?
                DEFAULT_STENCIL_TIMEOUT_MS : config.get("STENCIL_TIMEOUT_MS"));
        int backoffMs = Integer.parseInt(StringUtils.isBlank(config.get("STENCIL_BACKOFF_MS")) ?
                DEFAULT_STENCIL_BACKOFF_MS : config.get("STENCIL_BACKOFF_MS"));
        int retries = Integer.parseInt(StringUtils.isBlank(config.get("STENCIL_RETRIES")) ?
                DEFAULT_STENCIL_RETRIES : config.get("STENCIL_RETRIES"));
        int retryCount = retries;
        //get schema from server
        while (true) {
            try {
                log.info("fetching descriptors from {} with timeout: {}ms, backoff: {}ms {} retries pending", url, timeout, backoffMs, retryCount);
                descriptorBin = new RemoteFile().fetch(url, timeout);
                log.info("successfully fetched {}", url);
                break;
            } catch (IOException | RuntimeException e) {
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

    private void generateMap() {
        try {
            descriptorMap = new DescriptorMapBuilder().buildFrom(new ByteArrayInputStream(descriptorBin));
        } catch (IOException | Descriptors.DescriptorValidationException e) {
            throw new StencilRuntimeException(e);
        }
    }
}
