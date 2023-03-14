package com.gotocompany.stencil.http;

import java.io.IOException;

public interface RemoteFile {
    byte[] fetch(String url) throws IOException;
    void close() throws IOException;
}

