package api_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/stencil/server/models"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/status"
)

var (
	formError   = models.ErrMissingFormData.Message()
	uploadError = models.ErrUploadFailed.Message()
	success     = "success"
)

func createMultipartBody(name string, version string, dryrun bool) (*bytes.Buffer, *multipart.Writer, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	defer writer.Close()
	if err := writer.WriteField("name", name); err != nil {
		return nil, writer, err
	}
	if err := writer.WriteField("version", version); err != nil {
		return nil, writer, err
	}
	if err := writer.WriteField("dryrun", strconv.FormatBool(dryrun)); err != nil {
		return nil, writer, err
	}
	fileWriter, err := writer.CreateFormFile("file", "test.desc")
	if err != nil {
		return nil, writer, err
	}
	file, err := os.Open("./testdata/test.desc")
	if err != nil {
		return nil, writer, err
	}
	defer file.Close()
	_, err = io.Copy(fileWriter, file)
	return buf, writer, err
}

func TestUpload(t *testing.T) {
	for _, test := range []struct {
		desc         string
		name         string
		version      string
		exists       bool
		validateErr  error
		insertErr    error
		expectedCode int
		responseMsg  string
	}{
		{"should return 400 if name is missing", "", "1.0.1", false, nil, nil, 400, formError},
		{"should return 400 if version is missing", "name1", "", false, nil, nil, 400, formError},
		{"should return 400 if version is invalid semantic version", "name1", "invalid", false, nil, nil, 400, formError},
		{"should return 400 if backward check fails", "name1", "1.0.1", false, errors.New("validation"), nil, 400, "validation"},
		{"should return 409 if resource already exists", "name1", "1.0.1", true, nil, nil, 409, "Resource already exists"},
		{"should return 500 if insert fails", "name1", "1.0.1", false, nil, errors.New("insert fail"), 500, "Internal error"},
		{"should return 200 if upload succeeded", "name1", "1.0.1", false, nil, nil, 200, success},
	} {
		t.Run(test.desc, func(t *testing.T) {
			router, mockService, metadata, _ := setup()
			metadata.On("Exists", mock.Anything, mock.Anything).Return(test.exists)
			mockService.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(test.validateErr)
			mockService.On("Insert", mock.Anything, mock.Anything, mock.Anything).Return(test.insertErr)
			w := httptest.NewRecorder()
			body, writer, _ := createMultipartBody(test.name, test.version, false)
			req, _ := http.NewRequest("POST", "/v1/namespaces/namespace/descriptors", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)

			if w.Code == 200 {
				assert.JSONEq(t, fmt.Sprintf(`{"message": "%s", "dryrun": false}`, test.responseMsg), w.Body.String())
			} else {
				assert.JSONEq(t, fmt.Sprintf(`{"message": "%s"}`, test.responseMsg), w.Body.String())
			}
		})
		t.Run(fmt.Sprintf("gRPC: %s", test.desc), func(t *testing.T) {
			_, mockService, metadata, api := setup()
			metadata.On("Exists", mock.Anything, mock.Anything).Return(test.exists)
			mockService.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(test.validateErr)
			mockService.On("Insert", mock.Anything, mock.Anything, mock.Anything).Return(test.insertErr)
			data, err := os.ReadFile("./testdata/test.desc")
			assert.Nil(t, err)
			req := &stencilv1.UploadDescriptorRequest{
				Namespace: "namespace", Name: test.name, Version: test.version, Checks: &stencilv1.Checks{},
				Data: data,
			}
			res, err := api.UploadDescriptor(context.Background(), req)
			if test.expectedCode != 200 {
				e := status.Convert(err)
				assert.Equal(t, test.expectedCode, runtime.HTTPStatusFromCode(e.Code()))
			} else {
				assert.Equal(t, res.Dryrun, false)
				assert.Equal(t, res.Success, true)
				assert.Equal(t, res.Errors, "")
			}
		})
	}

	t.Run("should not insert if dry run flag is enabled", func(t *testing.T) {
		name := "name1"
		version := "1.0.1"
		router, mockService, metadata, _ := setup()
		metadata.On("Exists", mock.Anything, mock.Anything).Return(false)
		mockService.On("Validate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockService.On("Insert", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		w := httptest.NewRecorder()
		body, writer, _ := createMultipartBody(name, version, true)
		req, _ := http.NewRequest("POST", "/v1/namespaces/namespace/descriptors", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `{"message": "success", "dryrun": true}`, w.Body.String())
		mockService.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything, mock.Anything)
	})
}
