package io.odpf.stencil.utils;

import com.google.protobuf.Descriptors;
import io.odpf.stencil.models.DescriptorAndTypeName;

import java.util.HashMap;
import java.util.Map;

/**
 * Utility method to parse the types and packages based on the descriptors.
 */
public class StencilUtils {

    /**
     * Gets a map of proto package name and typeName from the supplied map of model descriptors
     *
     * @param allDescriptors - Stencil modelled descriptors
     * @return - map of proto package and java type
     */
    public static Map<String, String> getTypeNameToPackageNameMap(final Map<String, DescriptorAndTypeName> allDescriptors) {
        Map<String, String> typeNameMap = new HashMap();
        allDescriptors.entrySet().stream().forEach((mapEntry) -> {
            DescriptorAndTypeName descriptorAndTypeName = (DescriptorAndTypeName) mapEntry.getValue();
            if (descriptorAndTypeName != null) {
                typeNameMap.put(descriptorAndTypeName.getTypeName(), mapEntry.getKey());
            }
        });
        return typeNameMap;
    }

    /**
     * Gets a map of type and the descriptor from the supplied map of model descriptors
     *
     * @param allDescriptors - Stencil modelled descriptors
     * @return - map of type and the respective protobuff descriptor
     */
    public static Map<String, Descriptors.Descriptor> getAllProtobufDescriptors(final Map<String, DescriptorAndTypeName> allDescriptors) {
        Map<String, Descriptors.Descriptor> descriptorMap = new HashMap();
        allDescriptors.entrySet().stream().forEach((mapEntry) -> {
            DescriptorAndTypeName descriptorAndTypeName = (DescriptorAndTypeName) mapEntry.getValue();
            if (descriptorAndTypeName != null) {
                descriptorMap.put(mapEntry.getKey(), descriptorAndTypeName.getDescriptor());
            }
        });
        return descriptorMap;
    }
}
