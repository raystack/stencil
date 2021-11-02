package api

import (
	"context"

	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/odpf/stencil/server/schema"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createSchemaRequestToSchema(in *stencilv1.CreateSchemaRequest) *schema.Schema {
	return &schema.Schema{
		ID:          in.SchemaId,
		Data:        in.GetData(),
		NamespaceID: in.NamespaceId,
		Authority:   in.Authority,
		Format:      in.Format.String(),
		Description: in.Description,
	}
}

func schemaToSchemaProto(s *schema.Schema) *stencilv1.Schema {
	return &stencilv1.Schema{
		Id:          s.ID,
		Format:      stencilv1.Schema_Format(stencilv1.Schema_Format_value[s.Format]),
		Authority:   s.Authority,
		Description: s.Description,
		CreatedAt:   timestamppb.New(s.CreatedAt),
		UpdatedAt:   timestamppb.New(s.UpdatedAt),
	}
}

func (a *API) CreateSchema(ctx context.Context, in *stencilv1.CreateSchemaRequest) (*stencilv1.CreateSchemaResponse, error) {
	sc, err := a.SchemaService.CreateSchema(ctx, createSchemaRequestToSchema(in))
	return &stencilv1.CreateSchemaResponse{
		Schema: schemaToSchemaProto(sc),
	}, err
}

func (a *API) ListSchemas(ctx context.Context, in *stencilv1.ListSchemasRequest) (*stencilv1.ListSchemasResponse, error) {
	schemas, err := a.SchemaService.ListSchemas(ctx, in.Id)
	return &stencilv1.ListSchemasResponse{Schemas: schemas}, err
}
