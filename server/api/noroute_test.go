package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/odpf/stencil/server/config"
	server2 "github.com/odpf/stencil/server/server"

	"github.com/odpf/stencil/server/api"
	"github.com/stretchr/testify/assert"
)

func TestNoRoute(t *testing.T) {
	router := server2.Router(&api.API{}, &config.Config{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/random", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.JSONEq(t, `{"message": "page not found"}`, w.Body.String())
}
