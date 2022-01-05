package io.odpf.stencil.cache;

import java.io.IOException;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;

import com.google.protobuf.Descriptors;

import org.json.JSONArray;
import org.json.JSONObject;

import io.odpf.stencil.DescriptorMapBuilder;
import io.odpf.stencil.exception.StencilRuntimeException;
import io.odpf.stencil.http.RemoteFile;

public interface SchemaRefreshStrategy {
    Map<String, Descriptors.Descriptor> refresh(String url, RemoteFile remoteFile, final Map<String, Descriptors.Descriptor> prevDescriptor);

    static SchemaRefreshStrategy longPollingStrategy() {
        return (String url, RemoteFile remoteFile, Map<String, Descriptors.Descriptor> prevDescriptor) -> DescriptorMapBuilder.buildFrom(url,
                remoteFile);
    }

    static SchemaRefreshStrategy versionBasedRefresh() {
        final AtomicInteger lastVersion = new AtomicInteger();
        return (String url, RemoteFile remoteFile, Map<String, Descriptors.Descriptor> prevDescriptor) -> {
            try {
                byte[] data = remoteFile.fetch(String.format("%s/versions", url));
                JSONObject json = new JSONObject(new String(data));
                JSONArray versions = json.getJSONArray("versions");
                Integer maxVersion = 0;
                for (int i = 0; i < versions.length(); i++) {
                    if (versions.getInt(i) > maxVersion) {
                        maxVersion = versions.getInt(i);
                    }
                }
                if (maxVersion != 0 && maxVersion != lastVersion.get()) {
                    String newURL = String.format("%s/versions/%d", url, maxVersion);
                    Map<String, Descriptors.Descriptor> newSchema = DescriptorMapBuilder.buildFrom(newURL, remoteFile);
                    lastVersion.set(maxVersion);
                    return newSchema;
                }
                return prevDescriptor;
            } catch (IOException e) {
                throw new StencilRuntimeException(e);
            }
        };
    }
}
