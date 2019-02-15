package com.gojek.de.stencil;

import com.google.protobuf.Descriptors;
import org.junit.Test;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Map;
import java.util.Objects;

import static org.junit.Assert.assertNotNull;

public class DescriptorMapBuilderTest {
    @Test
    public void testStore() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, Descriptors.Descriptor> descriptorMap = new DescriptorMapBuilder().buildFrom(fileInputStream);
        assertNotNull(descriptorMap);
        assertNotNull(descriptorMap.get("com.gojek.stencil.TestMessage"));
    }
}
