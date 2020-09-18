package com.gojek.de.stencil.client;

import com.gojek.de.stencil.cache.DescriptorCacheLoader;
import com.gojek.de.stencil.config.StencilConfig;
import com.gojek.de.stencil.models.DescriptorAndTypeName;
import com.google.protobuf.Descriptors;

import java.io.IOException;
import java.io.Serializable;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;

public class MultiURLStencilClient implements Serializable, StencilClient {

    private List<StencilClient> stencilClients;

    public MultiURLStencilClient(List<String> urls, StencilConfig config, DescriptorCacheLoader cacheLoader) {
        stencilClients = urls.stream().map(url -> new URLStencilClient(url, config, cacheLoader)).collect(Collectors.toList());
    }

    @Override
    public Descriptors.Descriptor get(String protoClassName) {
        Optional<StencilClient> requiredStencil = stencilClients.stream().filter(stencilClient -> stencilClient.get(protoClassName) != null).findFirst();
        return requiredStencil.map(stencilClient -> stencilClient.get(protoClassName)).orElse(null);
    }

    @Override
    public Map<String, Descriptors.Descriptor> getAll() {
        Map<String, Descriptors.Descriptor> requiredStencil = new HashMap<>();
        stencilClients.stream().map(StencilClient::getAll)
                .forEach(requiredStencil::putAll);
        return requiredStencil;
    }

    @Override
    public Map<String, String> getTypeNameToPackageNameMap() {
        Map<String, String> requiredStencil = new HashMap<>();
        stencilClients.stream().map(StencilClient::getTypeNameToPackageNameMap)
                .forEach(requiredStencil::putAll);
        return requiredStencil;
    }

    @Override
    public Map<String, DescriptorAndTypeName> getAllDescriptorAndTypeName() {
        Map<String, DescriptorAndTypeName> requiredStencil = new HashMap<>();
        stencilClients.stream().map(StencilClient::getAllDescriptorAndTypeName)
                .forEach(requiredStencil::putAll);
        return requiredStencil;
    }

    @Override
    public void close() {
        stencilClients.forEach(c -> {
            try {
                c.close();
            } catch (IOException e) {
                e.printStackTrace();
            }
        });
    }

    @Override
    public void refresh() {
        stencilClients.forEach(c -> {
            c.refresh();
        });
    }
}
