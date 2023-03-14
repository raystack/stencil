package api_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goto/stencil/core/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHTTPGetSchema(t *testing.T) {
	nsName := "namespace1"
	schemaName := "scName"
	t.Run("should validate version number", func(t *testing.T) {
		_, _, _, mux, _ := setup()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1beta1/namespaces/%s/schemas/%s/versions/invalidNumber", nsName, schemaName), nil)
		mux.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
		assert.JSONEq(t, `{"code":2,"message":"invalid version number","details":[]}`, w.Body.String())
	})
	t.Run("should return http error if getSchema fails", func(t *testing.T) {
		version := int32(2)
		_, schemaSvc, _, mux, _ := setup()
		schemaSvc.On("Get", mock.Anything, nsName, schemaName, version).Return(nil, nil, errors.New("get error"))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1beta1/namespaces/%s/schemas/%s/versions/%d", nsName, schemaName, version), nil)
		mux.ServeHTTP(w, req)
		assert.Equal(t, 500, w.Code)
		assert.JSONEq(t, `{"code":2,"message":"get error","details":[]}`, w.Body.String())
	})
	t.Run("should return octet-stream content type for protobuf schema", func(t *testing.T) {
		version := int32(2)
		data := []byte("test data")
		_, schemaSvc, _, mux, _ := setup()
		schemaSvc.On("Get", mock.Anything, nsName, schemaName, version).Return(&schema.Metadata{Format: "FORMAT_PROTOBUF"}, data, nil)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1beta1/namespaces/%s/schemas/%s/versions/%d", nsName, schemaName, version), nil)
		mux.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, data, w.Body.Bytes())
		assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
	})
}

func TestHTTPSchemaCreate(t *testing.T) {
	nsName := "namespace"
	scName := "schemaName"
	format := "PROTOBUF"
	compatibility := "FULL"
	body := []byte("protobuf contents")
	t.Run("should return error if schema create fails", func(t *testing.T) {
		_, schemaSvc, _, mux, _ := setup()
		schemaSvc.On("Create", mock.Anything, nsName, scName, &schema.Metadata{Format: format, Compatibility: compatibility}, body).Return(schema.SchemaInfo{}, errors.New("create error"))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/v1beta1/namespaces/%s/schemas/%s", nsName, scName), bytes.NewBuffer(body))
		req.Header.Add("X-Format", format)
		req.Header.Add("X-Compatibility", compatibility)
		mux.ServeHTTP(w, req)
		assert.Equal(t, 500, w.Code)
		schemaSvc.AssertExpectations(t)
	})
	t.Run("should return schemaInfo in JSON after create", func(t *testing.T) {
		_, schemaSvc, _, mux, _ := setup()
		scInfo := schema.SchemaInfo{ID: "someID", Version: int32(2)}
		schemaSvc.On("Create", mock.Anything, nsName, scName, &schema.Metadata{Format: format, Compatibility: compatibility}, body).Return(scInfo, nil)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/v1beta1/namespaces/%s/schemas/%s", nsName, scName), bytes.NewBuffer(body))
		req.Header.Add("X-Format", format)
		req.Header.Add("X-Compatibility", compatibility)
		mux.ServeHTTP(w, req)
		assert.Equal(t, 201, w.Code)
		assert.JSONEq(t, `{"id": "someID", "location": "", "version": 2}`, w.Body.String())
		schemaSvc.AssertExpectations(t)
	})
}
