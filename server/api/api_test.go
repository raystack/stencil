package api_test

import (
	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/api/mocks"
)

func setup() (*mocks.NamespaceService, *mocks.SchemaService, *api.API) {
	nsService := &mocks.NamespaceService{}
	schemaService := &mocks.SchemaService{}
	v1 := &api.API{
		Namespace: nsService,
		Schema:    schemaService,
	}
	return nsService, schemaService, v1
}
