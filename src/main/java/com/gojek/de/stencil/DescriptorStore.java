package com.gojek.de.stencil;

import com.google.protobuf.Descriptors;

import java.util.HashMap;
import java.util.Map;

public class DescriptorStore {
    private static Map<String, String> configs = new HashMap<>();
    private static StencilClient client;

    public static void loadClientIfNull(Map<String, String> configs) {
        if (client == null) {
            DescriptorStore.configs = configs;
            load();
        }
    }

    private static void load() {
        if ("true".equals(configs.getOrDefault("ENABLE_STENCIL_URL", "false"))) {
            client = StencilClientFactory.getClient(
                    configs.getOrDefault("STENCIL_URL", ""),
                    configs
            );
        } else {
            client = StencilClientFactory.getClient();
        }
        client.load();
    }

    public static Descriptors.Descriptor get(String className) {
        if (client != null) {
            return client.get(className);
        } else {
            throw new StencilConfigurationException("DescriptorStore not loaded");
        }
    }
}
