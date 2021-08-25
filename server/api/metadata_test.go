package api_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/odpf/stencil/server/models"
	"github.com/odpf/stencil/server/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
			router, _, mockService, _ := setup()
			mockService.On("ListNames", mock.Anything, "namespace").Return(test.values, test.err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/v1/namespaces/namespace/descriptors", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.JSONEq(t, test.expectedResp, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}

}

// func TestListVersions(t *testing.T) {
// 	for _, test := range []struct {
// 		desc         string
// 		err          error
// 		values       []string
// 		expectedCode int
// 		expectedResp string
// 	}{
// 		{"should return list", nil, []string{"n1", "n2"}, 200, `["n1", "n2"]`},
// 		{"should return 404 if path not found", models.ErrNotFound, []string{}, 404, `{"message": "Not found"}`},
// 	} {
// 		t.Run(test.desc, func(t *testing.T) {
// 			router, _, mockService, _ := setup()
// 			mockService.On("ListVersions", mock.Anything, "namespace", "example").Return(test.values, test.err)

// 			w := httptest.NewRecorder()
// 			req, _ := http.NewRequest("GET", "/v1/namespaces/namespace/descriptors/example/versions", nil)
// 			router.ServeHTTP(w, req)

// 			assert.Equal(t, test.expectedCode, w.Code)
// 			assert.JSONEq(t, test.expectedResp, w.Body.String())
// 			mockService.AssertExpectations(t)
// 		})
// 	}
// }

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
			router, _, mockService, _ := setup()
			mockService.On("GetSnapshot", mock.Anything, "namespace", test.name, "", true).Return(&snapshot.Snapshot{Version: test.latestVersion}, test.err)
			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/namespaces/namespace/metadata/%s", test.name), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == 200 {
				expectedData := []byte(fmt.Sprintf(`{"version":"%s"}`, test.latestVersion))
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
			router, _, mockService, _ := setup()
			mockService.On("UpdateLatestVersion", mock.Anything, mock.Anything).Return(test.err)
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
