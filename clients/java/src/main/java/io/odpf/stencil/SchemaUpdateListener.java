package io.odpf.stencil;

import java.util.Map;
import com.google.protobuf.Descriptors;

public interface SchemaUpdateListener {
    void onSchemaUpdate(final Map<String, Descriptors.Descriptor> newDescriptor);
}
