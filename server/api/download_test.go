package api_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/odpf/stencil/server/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var downloadFail = errors.New("download fail")

func TestDownload(t *testing.T) {
	for _, test := range []struct {
		desc         string
		name         string
		version      string
		notFoundErr  error
		downloadErr  error
		expectedCode int
	}{
		{"should return 400 if name is missing", "", "1.0.1", nil, nil, 400},
		{"should return 400 if version is invalid", "name1", "invalid", nil, nil, 400},
		{"should return 404 if version is not found", "name1", "3.3.1", snapshot.ErrNotFound, nil, 404},
		{"should return 500 if finding snapshot fails", "name1", "3.3.1", errors.New("get snapshot fail"), nil, 500},
		{"should return 500 if download fails", "name1", "1.0.1", nil, downloadFail, 500},
		{"should return 200 if download succeeded", "name1", "1.0.1", nil, nil, 200},
		{"should be able to download with latest version", "name1", "latest", nil, nil, 200},
	} {
		t.Run(test.desc, func(t *testing.T) {
			router, mockService, mockMetadata, _ := setup()

			fileData := []byte("File contents")
			mockMetadata.On("GetSnapshot", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&snapshot.Snapshot{}, test.notFoundErr)
			mockService.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(fileData, test.downloadErr)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/namespaces/namespace/descriptors/%s/versions/%s", test.name, test.version), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == 200 {
				expectedHeader := fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, test.version, test.version)
				assert.Equal(t, []byte("File contents"), w.Body.Bytes())
				assert.Equal(t, expectedHeader, w.Header().Get("Content-Disposition"))
			}
		})
	}
}
