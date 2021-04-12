package api_test

import (
	server2 "github.com/odpf/stencil/server/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/odpf/stencil/server/api"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	router := server2.Router(&api.API{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message": "pong"}`, w.Body.String())
}
