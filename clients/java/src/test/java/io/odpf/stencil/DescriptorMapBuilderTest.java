package org.raystack.stencil;

import com.google.protobuf.Descriptors;
import org.junit.Test;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Map;
import java.util.Objects;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class DescriptorMapBuilderTest {
    @Test
    public void testStore() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, Descriptors.Descriptor> descriptorMap = DescriptorMapBuilder.buildFrom(fileInputStream);
        assertNotNull(descriptorMap);
        assertNotNull(descriptorMap.get("org.raystack.stencil.TestMessage"));
    }

    @Test
    public void TestNestedDescriptors() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, Descriptors.Descriptor> descriptorMap = DescriptorMapBuilder.buildFrom(fileInputStream);

        final Descriptors.Descriptor account_db_accounts = descriptorMap.get("org.raystack.stencil.account_db_accounts");
        assertNotNull(account_db_accounts.findFieldByName("id"));
        final Descriptors.Descriptor ID = descriptorMap.get("org.raystack.stencil.account_db_accounts.ID");
        assertNotNull(ID.findFieldByName("data"));
        final Descriptors.Descriptor fullDocument = descriptorMap.get("org.raystack.stencil.account_db_accounts.FULLDOCUMENT");
        assertNotNull(fullDocument.findFieldByName("customerid"));
        final Descriptors.Descriptor accounts_item = descriptorMap.get("org.raystack.stencil.account_db_accounts.FULLDOCUMENT.ACCOUNTS_ITEM");
        assertNotNull(accounts_item.findFieldByName("monthlyaveragebalance"));
    }

    @Test
    public void TestDescriptorsWithRecursiveFields() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, Descriptors.Descriptor> descriptorMap = DescriptorMapBuilder.buildFrom(fileInputStream);

        final Descriptors.Descriptor RecursiveLogMessage = descriptorMap.get("org.raystack.stencil.RecursiveLogMessage");
        assertNotNull(RecursiveLogMessage.findFieldByName("id"));
        final Descriptors.Descriptor RECORD = descriptorMap.get("org.raystack.stencil.RecursiveLogMessage.RECORD");
        assertNotNull(RECORD.findFieldByName("id"));
        assertNotNull(RECORD.findFieldByName("record"));
    }

    @Test
    public void TestDescriptorsWithoutPackageName() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, Descriptors.Descriptor> descriptorMap = DescriptorMapBuilder.buildFrom(fileInputStream);

        final Descriptors.Descriptor protoWithoutPackage = descriptorMap.get("org.raystack.stencil.RootField");
        assertEquals(".RootField", String.format(".%s", protoWithoutPackage.getFullName()));
    }

    @Test
    public void TestDescriptorsByProtoFullName() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, Descriptors.Descriptor> descriptorMap = DescriptorMapBuilder.buildFrom(fileInputStream);

        final Descriptors.Descriptor protoWithoutJavaPackage = descriptorMap.get("org.raystack.stencil.ImplicitOuterClass");
        assertNotNull(protoWithoutJavaPackage);
        assertEquals("org.raystack.stencil.ImplicitOuterClass", protoWithoutJavaPackage.getFullName());
    }

    @Test
    public void TestDescriptorsByProtoFullNameOrJavaName() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, Descriptors.Descriptor> descriptorMap = DescriptorMapBuilder.buildFrom(fileInputStream);

        Descriptors.Descriptor protoWithoutPackage = descriptorMap.get("org.raystack.stencil.RootField");
        assertEquals("RootField", protoWithoutPackage.getFullName());
        Descriptors.Descriptor descriptorByProtoName = descriptorMap.get("RootField");
        assertNotNull(descriptorByProtoName);
        assertEquals(protoWithoutPackage, descriptorByProtoName);
    }

}
