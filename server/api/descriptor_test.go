package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/odpf/stencil/server/config"
	server2 "github.com/odpf/stencil/server/server"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/api/mocks"
	"github.com/odpf/stencil/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	formError   = models.ErrMissingFormData.Message()
	uploadError = models.ErrUploadFailed.Message()
	success     = "success"
)

func setup() (*gin.Engine, *mocks.StoreService, *api.API) {
	mockService := &mocks.StoreService{}
	v1 := &api.API{
		Store: mockService,
	}
	router := server2.Router(v1, &config.Config{})
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
	for _, test := range []struct {
		desc         string
		err          error
		values       []string
		expectedCode int
		expectedResp string
	}{
		{"should return list", nil, []string{"n1", "n2"}, 200, `["n1", "n2"]`},
		{"should return 404 if path not found", models.ErrNotFound, []string{}, 404, `{"message": "Not found"}`},
	} {
		t.Run(test.desc, func(t *testing.T) {
			router, mockService, _ := setup()
			mockService.On("ListNames", "namespace").Return(test.values, test.err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/v1/namespaces/namespace/descriptors", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.JSONEq(t, test.expectedResp, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}

}

func TestListVersions(t *testing.T) {
	for _, test := range []struct {
		desc         string
		err          error
		values       []string
		expectedCode int
		expectedResp string
	}{
		{"should return list", nil, []string{"n1", "n2"}, 200, `["n1", "n2"]`},
		{"should return 404 if path not found", models.ErrNotFound, []string{}, 404, `{"message": "Not found"}`},
	} {
		t.Run(test.desc, func(t *testing.T) {
			router, mockService, _ := setup()
			mockService.On("ListVersions", "namespace", "example").Return(test.values, test.err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/v1/namespaces/namespace/descriptors/example/versions", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.JSONEq(t, test.expectedResp, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
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
			req, _ := http.NewRequest("POST", "/v1/namespaces/namespace/descriptors", body)
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

func TestGetVersion(t *testing.T) {
	for _, test := range []struct {
		desc          string
		name          string
		latestVersion string
		err           error
		expectedCode  int
	}{
		{"should return 500 if fetch version fails", "name1", "1.0.1", errors.New("fetch fail"), 500},
		{"should return latest version number", "name1", "1.0.2", nil, 200},
	} {
		t.Run(test.desc, func(t *testing.T) {
			router, mockService, _ := setup()
			metadata := &models.MetadataFile{Version: test.latestVersion, Updated: "2021-03-15T10:51:15Z"}
			mockService.On("GetMetadata", mock.Anything, mock.Anything).Return(metadata, test.err)
			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/namespaces/namespace/metadata/%s", test.name), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == 200 {
				expectedData, _ := json.Marshal(metadata)
				assert.Equal(t, expectedData, w.Body.Bytes())
			}
		})
	}
}

func TestUpdateLatestVersion(t *testing.T) {
	for _, test := range []struct {
		desc         string
		name         string
		version      string
		err          error
		expectedCode int
	}{
		{"should return 400 if name is missing", "", "1.0.1", nil, 400},
		{"should return 400 if version is missing", "name1", "", nil, 400},
		{"should return 400 if version not follows semantic verioning", "name1", "invalid0.1.0", nil, 400},
		{"should return 500 if store fails", "name1", "1.0.1", errors.New("store fail"), 500},
		{"should return success if update succeeds", "name1", "1.0.2", nil, 200},
	} {
		t.Run(test.desc, func(t *testing.T) {
			router, mockService, _ := setup()
			mockService.On("StoreMetadata", mock.Anything, mock.Anything).Return(test.err)
			w := httptest.NewRecorder()

			body := bytes.NewReader([]byte(fmt.Sprintf(`{"name": "%s", "version": "%s"}`, test.name, test.version)))
			req, _ := http.NewRequest("POST", "/v1/namespaces/namespace/metadata", body)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == 200 {
				assert.JSONEq(t, `{"message": "success"}`, w.Body.String())
				mockService.AssertExpectations(t)
			}
		})
	}
}
