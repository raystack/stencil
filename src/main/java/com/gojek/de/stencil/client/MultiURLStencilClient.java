package com.gojek.de.stencil.client;

import com.google.common.cache.CacheLoader;
import com.google.protobuf.Descriptors;

import java.io.Serializable;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;

public class MultiURLStencilClient implements Serializable, StencilClient {


    private List<StencilClient> stencilClients;

    public MultiURLStencilClient(List<String> urls, Map<String, String> config, CacheLoader cacheLoader) {
        stencilClients = urls.stream().map(url -> new URLStencilClient(url, config, cacheLoader)).collect(Collectors.toList());
    }

    @Override
    public Descriptors.Descriptor get(String protoClassName) {
        Optional<StencilClient> requiredStencil = stencilClients.stream().filter(stencilClient -> stencilClient.get(protoClassName) != null).findFirst();
        return requiredStencil.map(stencilClient -> stencilClient.get(protoClassName)).orElse(null);
    }
}
