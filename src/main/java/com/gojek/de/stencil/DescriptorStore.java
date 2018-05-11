package com.gojek.de.stencil;

import com.google.protobuf.Descriptors;

import java.util.HashMap;
import java.util.Map;

public class DescriptorStore {
    private static Map<String, String> configs = new HashMap<>();
    private static Map<String, Descriptors.Descriptor> descriptorMap = new HashMap<>();
    private static StencilClient client;

    public static void setConfigs(Map<String, String> configs) {
        DescriptorStore.configs = configs;
    }

    public static void load() {
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
        return client.get(className);
    }
}
