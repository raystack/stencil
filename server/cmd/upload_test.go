package cmd

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/odpf/stencil/server/config"
	"github.com/stretchr/testify/assert"
)

func TestUploadDescriptor(t *testing.T) {
	t.Run("should successfully upload descriptor", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.ParseMultipartForm(50 * 1024 * 1024)
			path := r.URL.Path
			name := r.Form.Get("name")
			version := r.Form.Get("version")
			descFile, _, _ := r.FormFile("file")
			defer descFile.Close()
			h := md5.New()
			if _, err := io.Copy(h, descFile); err != nil {
				t.Fatalf("Error reading file from request: %v", err)
			}

			fileDigest := fmt.Sprintf("%x", h.Sum(nil))

			assert.Equal(t, "b05403212c66bdc8ccc597fedf6cd5fe", fileDigest)
			assert.Equal(t, "test-entities", name)
			assert.Equal(t, "0.0.1", version)
			assert.Equal(t, "/v1/namespaces/pilot/descriptors", path)
			w.WriteHeader(200)
		}))
		defer ts.Close()

		testURL, _ := url.Parse(ts.URL)
		config := config.LoadConfig()
		config.Scheme = testURL.Scheme
		config.Host = strings.Split(testURL.Host, ":")[0]
		config.Port = strings.Split(testURL.Host, ":")[1]

		err := uploadDescriptor(config, "pilot", "test-entities", "0.0.1", "./test_data/test_file")
		if err != nil {
			t.Fatal(err)
		}
	})
}
