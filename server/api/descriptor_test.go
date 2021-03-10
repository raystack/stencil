package api_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server"
	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/api/mocks"
	"github.com/odpf/stencil/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	formError     = models.ErrMissingFormData.Message
	uploadError   = models.ErrUploadFailed.Message
	downloadError = models.ErrDownloadFailed.Message
	success       = "success"
)

func setup() (*gin.Engine, *mocks.StoreService, *api.API) {
	mockService := &mocks.StoreService{}
	v1 := &api.API{
		Store: mockService,
	}
	router := server.Router(v1)
	return router, mockService, v1
}

func createMultipartBody(name string, version string) (*bytes.Buffer, *multipart.Writer, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	defer writer.Close()
	if err := writer.WriteField("name", name); err != nil {
		return nil, writer, err
	}
	if err := writer.WriteField("version", version); err != nil {
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

func mockFileData(contents string) *models.FileData {
	reader := bytes.NewReader([]byte(contents))
	r := ioutil.NopCloser(reader)
	fileData := &models.FileData{
		Reader:        r,
		ContentLength: reader.Size(),
	}
	return fileData
}

func TestList(t *testing.T) {
	t.Run("should return list", func(t *testing.T) {
		router, mockService, _ := setup()
		mockService.On("ListNames", "org").Return([]string{"name1", "name2"})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/descriptors", nil)
		req.Header.Set("x-scope-orgid", "org")
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `["name1", "name2"]`, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("should return 400 if org id not specified", func(t *testing.T) {
		router, _, _ := setup()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/descriptors", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.JSONEq(t, `{"message": "x-scope-orgid header should be present"}`, w.Body.String())
	})

}

func TestListVersions(t *testing.T) {
	t.Run("should return list", func(t *testing.T) {
		router, mockService, _ := setup()
		mockService.On("ListVersions", "org", "example").Return([]string{"name1", "name2"})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/descriptors/example", nil)
		req.Header.Set("x-scope-orgid", "org")
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `["name1", "name2"]`, w.Body.String())
		mockService.AssertExpectations(t)
	})
}

func TestUpload(t *testing.T) {
	for _, test := range []struct {
		desc         string
		name         string
		version      string
		uploadErr    error
		expectedCode int
		responseMsg  string
	}{
		{"should return 400 if name is missing", "", "1.0.1", nil, 400, formError},
		{"should return 400 if version is missing", "name1", "", nil, 400, formError},
		{"should return 400 if version is invalid semantic version", "name1", "invalid", nil, 400, formError},
		{"should return 500 if upload fails", "name1", "1.0.1", errors.New("upload fail"), 500, uploadError},
		{"should return 200 if upload succeeded", "name1", "1.0.1", nil, 200, success},
	} {
		t.Run(test.desc, func(t *testing.T) {
			router, mockService, _ := setup()
			mockService.On("Upload", mock.Anything, mock.Anything).Return(test.uploadErr)
			w := httptest.NewRecorder()
			body, writer, _ := createMultipartBody(test.name, test.version)
			req, _ := http.NewRequest("POST", "/v1/descriptors", body)
			req.Header.Set("x-scope-orgid", "org")
			req.Header.Set("Content-Type", writer.FormDataContentType())

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.JSONEq(t, fmt.Sprintf(`{"message": "%s"}`, test.responseMsg), w.Body.String())
		})
	}
}

func TestDownload(t *testing.T) {
	for _, test := range []struct {
		desc         string
		name         string
		version      string
		downloadErr  error
		expectedCode int
	}{
		{"should return 400 if name is missing", "", "1.0.1", nil, 400},
		{"should return 400 if version is invalid", "name1", "invalid", nil, 400},
		{"should return 500 if download fails", "name1", "1.0.1", errors.New("download fail"), 500},
		{"should return 200 if download succeeded", "name1", "1.0.1", nil, 200},
		{"should be able to download with latest version", "name1", "latest", nil, 200},
	} {
		t.Run(test.desc, func(t *testing.T) {
			router, mockService, _ := setup()

			fileData := mockFileData("File contents")
			mockService.On("Download", mock.Anything, mock.Anything).Return(fileData, test.downloadErr)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/descriptors/%s/%s", test.name, test.version), nil)
			req.Header.Set("x-scope-orgid", "org")

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
