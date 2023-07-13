package org.raystack.stencil.client;

import com.google.protobuf.Descriptors;
import org.raystack.stencil.cache.SchemaCacheLoader;
import org.raystack.stencil.config.StencilConfig;

import java.io.IOException;
import java.io.Serializable;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;

/**
* {@link StencilClient} implementation that can fetch descriptor sets from multiple URLs
*/
public class MultiURLStencilClient implements Serializable, StencilClient {

    private List<StencilClient> stencilClients;

    /**
     * @param urls List of URLs to fetch protobuf descriptor sets from
     * @param config Stencil configs
     * @param cacheLoader Extension of Guava {@link com.google.common.cache.CacheLoader} for Proto Descriptor sets
     */
    public MultiURLStencilClient(List<String> urls, StencilConfig config, SchemaCacheLoader cacheLoader) {
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
