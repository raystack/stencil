package io.odpf.stencil.client;

import io.odpf.stencil.StencilClientFactory;
import com.google.protobuf.Descriptors;
import org.junit.Test;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.ObjectOutputStream;

import static org.junit.Assert.*;

public class ClassLoadStencilClientTest {

    private static final String LOOKUP_KEY = "io.odpf.stencil.TestMessage";

    @Test
    public void getDescriptorFromClassPath() {
        StencilClient c = StencilClientFactory.getClient();
        Descriptors.Descriptor desc = c.get(LOOKUP_KEY);
        assertNotNull(desc);
    }

    @Test
    public void ClassNotPresent() {
        StencilClient c = StencilClientFactory.getClient();
        Descriptors.Descriptor dsc = c.get("non_existent_proto");
        assertNull(dsc);
    }

    @Test
    public void shouldBeSerializable() throws IOException {
        serializeObject(StencilClientFactory.getClient());
    }

    public static byte[] serializeObject(Object o) throws IOException {
        try (ByteArrayOutputStream baos = new ByteArrayOutputStream();
             ObjectOutputStream oos = new ObjectOutputStream(baos)) {
            oos.writeObject(o);
            oos.flush();
            return baos.toByteArray();
        }
    }
}
