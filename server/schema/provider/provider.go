package provider

import "github.com/odpf/stencil/server/schema"

type schemaFun func([]byte) (*schema.SchemaFile, error)

type SchemaProvider struct {
	mapper map[string]schemaFun
}

func (s *SchemaProvider) GetSchemaFile(format string, data []byte) (*schema.SchemaFile, error) {
	return pbGetSchemaFile(data)
}

func NewSchemaProvider() *SchemaProvider {
	mp := make(map[string]schemaFun)
	mp["PROTOBUF"] = pbGetSchemaFile
	return &SchemaProvider{
		mapper: mp,
	}
}
