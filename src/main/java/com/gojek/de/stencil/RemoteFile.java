package com.gojek.de.stencil;

import java.io.IOException;

public interface RemoteFile {
    byte[] fetch(String url, int timeout) throws IOException;
}

