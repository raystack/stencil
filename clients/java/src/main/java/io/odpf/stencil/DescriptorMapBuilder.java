package io.odpf.stencil;

import com.google.protobuf.DescriptorProtos;
import com.google.protobuf.Descriptors;

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

    public Map<String, Descriptors.Descriptor> buildFrom(InputStream stream) throws IOException, Descriptors.DescriptorValidationException {
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
            fd.getMessageTypes().stream().forEach(desc -> descriptorMap.putAll(getFlattenedDescriptors(desc, javaPackage, protoPackage, "", new HashMap<>())));
        });

        return descriptorMap;
    }

    private Map<String, Descriptors.Descriptor> getFlattenedDescriptors(Descriptors.Descriptor descriptor, String javaPackage, String protoPackage, String parentClassName, Map<String, Descriptors.Descriptor> initialDescriptorMap) {
        String className = getClassName(descriptor, parentClassName);
        String javaClassName = javaPackage.isEmpty() ? className : String.format("%s.%s", javaPackage, className);
        initialDescriptorMap.put(javaClassName, descriptor);
        descriptor.getNestedTypes()
                .forEach(desc -> getFlattenedDescriptors(desc, javaPackage, protoPackage, className, initialDescriptorMap));
        return initialDescriptorMap;
    }


    private String getClassName(Descriptors.Descriptor descriptor, String parentClassName) {
        return parentClassName.isEmpty() ? descriptor.getName() : parentClassName + "." + descriptor.getName();
    }

}
