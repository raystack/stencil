package provider

import (
	"errors"

	"github.com/odpf/stencil/server/avro"
	"github.com/odpf/stencil/server/json"
	"github.com/odpf/stencil/server/protobuf"
	"github.com/odpf/stencil/server/schema"
)

type parseFn func([]byte) (schema.ParsedSchema, error)

type SchemaProvider struct {
	mapper map[string]parseFn
}

func (s *SchemaProvider) ParseSchema(format string, data []byte) (schema.ParsedSchema, error) {
	fn, ok := s.mapper[format]
	if ok {
		return fn(data)
	}
	return nil, errors.New("unknown schema")
}

func NewSchemaProvider() *SchemaProvider {
	mp := make(map[string]parseFn)
	mp["FORMAT_PROTOBUF"] = protobuf.GetParsedSchema
	mp["FORMAT_AVRO"] = avro.ParseSchema
	mp["FORMAT_JSON"] = json.GetParsedSchema
	return &SchemaProvider{
		mapper: mp,
	}
}
