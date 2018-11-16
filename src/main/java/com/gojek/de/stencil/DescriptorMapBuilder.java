package com.gojek.de.stencil;

import com.google.protobuf.DescriptorProtos;
import com.google.protobuf.Descriptors;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;

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
            fd.getMessageTypes().stream().forEach(desc -> {
                String className = desc.getName();
                descriptorMap.put(String.format("%s.%s", javaPackage, className), desc);
            });
        });

        return descriptorMap;
    }
}
