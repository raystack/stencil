package api_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/stencil/models"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/status"
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
		{"should return 404 if version is not found", "name1", "3.3.1", models.ErrSnapshotNotFound, nil, 404},
		{"should return 500 if finding snapshot fails", "name1", "3.3.1", errors.New("get snapshot fail"), nil, 500},
		{"should return 500 if download fails", "name1", "1.0.1", nil, downloadFail, 500},
		{"should return 200 if download succeeded", "name1", "1.0.1", nil, nil, 200},
		{"should be able to download with latest version", "name1", "latest", nil, nil, 200},
	} {
		t.Run(fmt.Sprintf("http: %s", test.desc), func(t *testing.T) {
			router, mockService, mockMetadata, _ := setup()

			fileData := []byte("File contents")
			mockMetadata.On("GetSnapshotByFields", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&models.Snapshot{}, test.notFoundErr)
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
		t.Run(fmt.Sprintf("gRPC: %s", test.desc), func(t *testing.T) {
			ctx := context.Background()
			_, mockService, mockMetadata, a := setup()

			fileData := []byte("File contents")
			mockMetadata.On("GetSnapshotByFields", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&models.Snapshot{}, test.notFoundErr)
			mockService.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(fileData, test.downloadErr)
			req := &stencilv1.DownloadDescriptorRequest{Namespace: "namespace", Name: test.name, Version: test.version}
			res, err := a.DownloadDescriptor(ctx, req)
			if test.expectedCode != 200 {
				e := status.Convert(err)
				assert.Equal(t, test.expectedCode, runtime.HTTPStatusFromCode(e.Code()))
			} else {
				assert.Equal(t, res.Data, []byte("File contents"))
			}
		})
	}
	t.Run("should return 404 if file content not found", func(t *testing.T) {
		router, mockService, mockMetadata, _ := setup()
		fileData := []byte("")
		mockMetadata.On("GetSnapshotByFields", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&models.Snapshot{}, nil)
		mockService.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(fileData, nil)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/namespaces/namespace/descriptors/n/versions/latest", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Equal(t, []byte(`{"message":"not found"}`), w.Body.Bytes())
	})
	t.Run("should return schema for given type", func(t *testing.T) {

		router, mockService, mockMetadata, _ := setup()

		fileData := []byte("File contents")
		typeArg := "odpf.stencil.v1.Snapshot"
		mockMetadata.On("GetSnapshotByFields", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&models.Snapshot{}, nil)
		mockService.On("Get", mock.Anything, mock.Anything, []string{typeArg}).Return(fileData, nil)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/namespaces/namespace/descriptors/name/versions/0.0.1/types/%s", typeArg), nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, []byte("File contents"), w.Body.Bytes())
		mockService.AssertExpectations(t)
	})
	t.Run("should return schema for given multiple types", func(t *testing.T) {

		router, mockService, mockMetadata, _ := setup()

		fileData := []byte("File contents")
		typesArg := []string{"odpf.stencil.v1.Snapshot", "odpf.stencil.v1.Checks"}
		mockMetadata.On("GetSnapshotByFields", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&models.Snapshot{}, nil)
		mockService.On("Get", mock.Anything, mock.Anything, typesArg).Return(fileData, nil)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/namespaces/namespace/descriptors/name/versions/0.0.1?fullnames=%s&fullnames=%s", typesArg[0], typesArg[1]), nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, []byte("File contents"), w.Body.Bytes())
		mockService.AssertExpectations(t)
	})
}
