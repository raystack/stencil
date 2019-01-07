package com.gojek.de.stencil;

import com.google.protobuf.Descriptors;

import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;

public class MultiURLStencilClient extends StencilClient {


    private List<StencilClient> stencilClients;

    public MultiURLStencilClient(List<String> urls, Map<String, String> config) {
        stencilClients = urls.stream().map(url -> new URLStencilClient(url, config)).collect(Collectors.toList());
    }

    @Override
    public Descriptors.Descriptor get(String protoClassName) {
        Optional<StencilClient> requiredStencil = stencilClients.stream().filter(stencilClient -> stencilClient.get(protoClassName) != null).findFirst();
        return requiredStencil.map(stencilClient -> stencilClient.get(protoClassName)).orElse(null);
    }
}
