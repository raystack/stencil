package io.odpf.stencil;

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
        Map<String, Descriptors.Descriptor> descriptorMap = new DescriptorMapBuilder().buildFrom(fileInputStream);
        assertNotNull(descriptorMap);
        assertNotNull(descriptorMap.get("io.odpf.stencil.TestMessage"));
    }

    @Test
    public void TestNestedDescriptors() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, Descriptors.Descriptor> descriptorMap = new DescriptorMapBuilder().buildFrom(fileInputStream);

        final Descriptors.Descriptor account_db_accounts = descriptorMap.get("io.odpf.stencil.account_db_accounts");
        assertNotNull(account_db_accounts.findFieldByName("id"));
        final Descriptors.Descriptor ID = descriptorMap.get("io.odpf.stencil.account_db_accounts.ID");
        assertNotNull(ID.findFieldByName("data"));
        final Descriptors.Descriptor fullDocument = descriptorMap.get("io.odpf.stencil.account_db_accounts.FULLDOCUMENT");
        assertNotNull(fullDocument.findFieldByName("customerid"));
        final Descriptors.Descriptor accounts_item = descriptorMap.get("io.odpf.stencil.account_db_accounts.FULLDOCUMENT.ACCOUNTS_ITEM");
        assertNotNull(accounts_item.findFieldByName("monthlyaveragebalance"));
    }

    @Test
    public void TestDescriptorsWithRecursiveFields() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, Descriptors.Descriptor> descriptorMap = new DescriptorMapBuilder().buildFrom(fileInputStream);

        final Descriptors.Descriptor RecursiveLogMessage = descriptorMap.get("io.odpf.stencil.RecursiveLogMessage");
        assertNotNull(RecursiveLogMessage.findFieldByName("id"));
        final Descriptors.Descriptor RECORD = descriptorMap.get("io.odpf.stencil.RecursiveLogMessage.RECORD");
        assertNotNull(RECORD.findFieldByName("id"));
        assertNotNull(RECORD.findFieldByName("record"));
    }

    @Test
    public void TestDescriptorsWithoutPackageName() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, Descriptors.Descriptor> descriptorMap = new DescriptorMapBuilder().buildFrom(fileInputStream);

        final Descriptors.Descriptor protoWithoutPackage = descriptorMap.get("io.odpf.stencil.RootField");
        assertEquals(".RootField", String.format(".%s", protoWithoutPackage.getFullName()));
    }

}
