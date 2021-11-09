package provider

import (
	"github.com/google/uuid"
	"github.com/odpf/stencil/server/schema"
)

func pbGetSchemaFile(data []byte) (*schema.SchemaFile, error) {
	id := uuid.NewSHA1(uuid.NameSpaceOID, data)
	return &schema.SchemaFile{
		ID:   id.String(),
		Data: data,
	}, nil
}
