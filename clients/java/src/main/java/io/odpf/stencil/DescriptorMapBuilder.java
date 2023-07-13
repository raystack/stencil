package org.raystack.stencil;

import com.google.protobuf.DescriptorProtos;
import com.google.protobuf.Descriptors;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.raystack.stencil.exception.StencilRuntimeException;
import org.raystack.stencil.http.RemoteFile;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;

/**
 * Provides methods to generate a map of Protobuf
 * descriptor and it's name for future lookups
 */
public class DescriptorMapBuilder {

    private static final Logger logger = LoggerFactory.getLogger(DescriptorMapBuilder.class);

    public static Map<String, Descriptors.Descriptor> buildFrom(String url, RemoteFile remoteFile) {
        try {
            logger.info("fetching descriptors from {}", url);
            byte[] descriptorBin = remoteFile.fetch(url);
            logger.info("successfully fetched {}", url);
            InputStream inputStream = new ByteArrayInputStream(descriptorBin);
            Map<String, Descriptors.Descriptor> newDescriptorsMap = DescriptorMapBuilder.buildFrom(inputStream);
            return newDescriptorsMap;
        } catch (IOException | Descriptors.DescriptorValidationException e) {
            throw new StencilRuntimeException(e);
        }
    }

    public static Map<String, Descriptors.Descriptor> buildFrom(InputStream stream) throws IOException, Descriptors.DescriptorValidationException {
        Map<String, Descriptors.Descriptor> descriptorMap = new HashMap<>();
        ArrayList<Descriptors.FileDescriptor> fileDescriptors = new ArrayList<>();
        DescriptorProtos.FileDescriptorSet descriptorSet = DescriptorProtos.FileDescriptorSet.parseFrom(stream);

        for (DescriptorProtos.FileDescriptorProto fdp : descriptorSet.getFileList()) {
            fileDescriptors.add(
                    Descriptors.FileDescriptor.buildFrom(fdp, fileDescriptors.toArray(new Descriptors.FileDescriptor[0]))
            );
        }

        fileDescriptors.forEach(fd -> {
            String javaPackage = fd.getOptions().getJavaPackage();
            String protoPackage = fd.getPackage();
            fd.getMessageTypes().stream().forEach(desc -> descriptorMap.putAll(getFlattenedDescriptors(desc, javaPackage, protoPackage, new HashMap<>())));
        });

        return descriptorMap;
    }

    private static Map<String, Descriptors.Descriptor> getFlattenedDescriptors(Descriptors.Descriptor descriptor, String javaPackage, String protoPackage, Map<String, Descriptors.Descriptor> initialDescriptorMap) {
        String fullName = descriptor.getFullName();
        initialDescriptorMap.put(fullName, descriptor);
        if (!javaPackage.isEmpty() && !javaPackage.equals(protoPackage)) {
            initialDescriptorMap.put(getClassName(descriptor, protoPackage, javaPackage), descriptor);
        }
        descriptor.getNestedTypes()
                .forEach(desc -> getFlattenedDescriptors(desc, javaPackage, protoPackage, initialDescriptorMap));
        return initialDescriptorMap;
    }

    private static String getClassName(Descriptors.Descriptor descriptor, String protoPackage, String javaPackage) {
        if (protoPackage.isEmpty()) {
            return String.format("%s.%s", javaPackage, descriptor.getFullName());
        }
        return descriptor.getFullName().replaceFirst(protoPackage, javaPackage);
    }
}
