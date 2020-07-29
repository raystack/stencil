package com.gojek.de.stencil;

import com.gojek.de.stencil.models.DescriptorAndTypeName;
import com.google.protobuf.DescriptorProtos;
import com.google.protobuf.Descriptors;
import org.junit.Test;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
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
        Map<String, DescriptorAndTypeName> descriptorMap = new DescriptorMapBuilder().buildFrom(fileInputStream);
        assertNotNull(descriptorMap);
        assertNotNull(descriptorMap.get("com.gojek.stencil.TestMessage"));
    }

    @Test
    public void TestNestedDescriptors() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, DescriptorAndTypeName> descriptorMap = new DescriptorMapBuilder().buildFrom(fileInputStream);

        final DescriptorAndTypeName account_db_accounts = descriptorMap.get("com.gojek.stencil.account_db_accounts");
        assertNotNull(account_db_accounts.getDescriptor().findFieldByName("id"));
        final DescriptorAndTypeName ID = descriptorMap.get("com.gojek.stencil.account_db_accounts.ID");
        assertNotNull(ID.getDescriptor().findFieldByName("data"));
        final DescriptorAndTypeName fullDocument = descriptorMap.get("com.gojek.stencil.account_db_accounts.FULLDOCUMENT");
        assertNotNull(fullDocument.getDescriptor().findFieldByName("customerid"));
        final DescriptorAndTypeName accounts_item = descriptorMap.get("com.gojek.stencil.account_db_accounts.FULLDOCUMENT.ACCOUNTS_ITEM");
        assertNotNull(accounts_item.getDescriptor().findFieldByName("monthlyaveragebalance"));
    }

    @Test
    public void TestDescriptorsWithRecursiveFields() throws IOException, Descriptors.DescriptorValidationException {
        ClassLoader classLoader = getClass().getClassLoader();
        String descriptorFilePath = "__files/descriptors.bin";
        InputStream fileInputStream = new FileInputStream(Objects.requireNonNull(classLoader.getResource(descriptorFilePath)).getFile());
        Map<String, DescriptorAndTypeName> descriptorMap = new DescriptorMapBuilder().buildFrom(fileInputStream);

        final DescriptorAndTypeName RecursiveLogMessage = descriptorMap.get("com.gojek.stencil.RecursiveLogMessage");
        assertNotNull(RecursiveLogMessage.getDescriptor().findFieldByName("id"));
        final DescriptorAndTypeName RECORD = descriptorMap.get("com.gojek.stencil.RecursiveLogMessage.RECORD");
        assertNotNull(RECORD.getDescriptor().findFieldByName("id"));
        assertNotNull(RECORD.getDescriptor().findFieldByName("record"));
    }

}
