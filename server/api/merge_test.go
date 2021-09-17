package api_test

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/odpf/stencil/models"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "io"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "testing"
)

func createMultipart() (*bytes.Buffer, *multipart.Writer, error) {
    buf := new(bytes.Buffer)
    writer := multipart.NewWriter(buf)
    defer writer.Close()
    fileWriter, err := writer.CreateFormFile("file", "test.desc")
    if err != nil {
        return nil, writer, err
    }
    _, err = io.Copy(fileWriter, bytes.NewReader([]byte("new data")))
    return buf, writer, err
}

func TestMerge(t *testing.T) {

    for _, test := range []struct {
        desc         string
        name         string
        version      string
        notFoundErr  error
        getPrevErr   error
        mergeErr     error
        expectedCode int
    }{
        {"should return 400 if name is missing", "", "1.0.1", nil, nil, nil, 400},
        {"should return 404 if version missing", "name1", "", nil, nil, nil, 404},
        {"should return 400 if version is invalid", "name1", "invalid", nil, nil, nil, 400},
        {"should return 404 if version is not found", "name1", "3.3.1", models.ErrSnapshotNotFound, nil, nil, 404},
        {"should return 500 if finding snapshot fails", "name1", "3.3.1", errors.New("get snapshot fail"), nil, nil, 500},
        {"should return 500 if getting previous data fails", "name1", "1.0.1", nil, errors.New("get prev data failed"), nil, 500},
        {"should return 500 if merge fails", "name1", "1.0.1", nil, nil, errors.New("merge failed"), 500},
        {"should return 200 if merge succeeded", "name1", "1.0.1", nil, nil, nil, 200},
        {"should be able to merge with latest version", "name1", "latest", nil, nil, nil, 200},
    } {
        t.Run(fmt.Sprintf("http: %s", test.desc), func(t *testing.T) {
            router, mockService, mockMetadata, _ := setup()
            prevData := []byte("prev data")
            mergedData := []byte("merged data")
            mockMetadata.On("GetSnapshotByFields", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&models.Snapshot{}, test.notFoundErr)
            mockService.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(prevData, test.getPrevErr)
            mockService.On("Merge", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mergedData, test.mergeErr)
            w := httptest.NewRecorder()
            body, writer, _ := createMultipart()
            req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/namespaces/namespace/descriptors/%s/versions/%s", test.name, test.version), body)
            req.Header.Set("Content-Type", writer.FormDataContentType())
            router.ServeHTTP(w, req)
            assert.Equal(t, test.expectedCode, w.Code)
            if test.expectedCode == 200 {
                expectedHeader := fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, test.name, test.name)
                assert.Equal(t, []byte("merged data"), w.Body.Bytes())
                assert.Equal(t, expectedHeader, w.Header().Get("Content-Disposition"))
            }
        })
    }

    t.Run("should return 400 if file body is missing", func(t *testing.T) {
        router, _, _, _ := setup()
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/namespaces/namespace/descriptors/%s/versions/%s", "name1", "3.3.1"), nil)
        router.ServeHTTP(w, req)
        assert.Equal(t, 400, w.Code)
    })

}