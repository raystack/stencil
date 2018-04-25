package com.gojek.esb.stencil;

import com.gojek.de.stencil.DescriptorStore;
import com.gojek.de.stencil.DescriptorStoreBuilder;
import com.google.protobuf.Descriptors;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;

import java.io.IOException;
import java.io.InputStream;

@Slf4j
public class StencilHttpClient {
    private Builder builder;
    private DescriptorStore descriptorStore;

    public Descriptors.Descriptor getForClass(String className) {
        return descriptorStore.getForClassName(className);
    }

    public void reload() {
        builder.build();
    }

    public class Builder {
        public Builder() {
            builder = this;
        }
        private String schemaDir;
        private String schemaVersion;
        private int retries = 3;
        private int backoffMs = 1000;
        private int timeoutMs = 10000;

        public void build() {
            if (StringUtils.isBlank(schemaDir)
                    || StringUtils.isBlank(schemaVersion)) {
                throw new StencilConfigurationException(
                        String.format("schemaDir: %s, SchemaVersion: %s cannot be blank")
                );
            }

            int retryCount = retries;
            //get schema from server
            while (true) {
                try {
                    InputStream is = new RemoteFile().fetch(getUrl(), timeoutMs);
                    descriptorStore = new DescriptorStoreBuilder().buildFrom(is);
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

        private String getUrl() {
            String url;
            if (schemaDir.endsWith("/")) {
                url = schemaDir;
            } else {
                url = schemaDir + "/";
            }
            return url + schemaVersion;
        }

        public Builder setSchemaDir(String schemaDir) {
            this.schemaDir = schemaDir;
            return this;
        }

        public Builder setSchemaVersion(String schemaVersion) {
            this.schemaVersion = schemaVersion;
            return this;
        }

        public Builder setRetries(int retries) {
            this.retries = retries;
            return this;
        }

        public Builder setBackoffMs(int backoffMs) {
            this.backoffMs = backoffMs;
            return this;
        }

        public Builder setTimeoutMs(int timeoutMs) {
            this.timeoutMs = timeoutMs;
            return this;
        }
    }
}
