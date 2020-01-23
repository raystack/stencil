package com.gojek.de.stencil;

import com.gojek.de.stencil.models.DescriptorAndTypeName;
import com.google.protobuf.DescriptorProtos;
import com.google.protobuf.Descriptors;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;

public class DescriptorMapBuilder {

    public Map<String, DescriptorAndTypeName> buildFrom(InputStream stream) throws IOException, Descriptors.DescriptorValidationException {
        Map<String, DescriptorAndTypeName> descriptorMap = new HashMap<>();
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
            fd.getMessageTypes().stream().forEach(desc -> {
                String className = desc.getName();
                desc.getNestedTypes().stream().forEach(nestedDesc -> {
                    String nestedClassName = nestedDesc.getName();
                    descriptorMap.put(
                            String.format("%s.%s.%s", javaPackage, className, nestedClassName),
                            new DescriptorAndTypeName(
                                    nestedDesc,
                                    String.format(".%s.%s.%s", protoPackage, className, nestedClassName)
                            ));
                });
                descriptorMap.put(
                        String.format("%s.%s", javaPackage, className),
                        new DescriptorAndTypeName(
                                desc,
                                String.format(".%s.%s", protoPackage, className)
                        ));
            });
        });

        return descriptorMap;
    }
}
