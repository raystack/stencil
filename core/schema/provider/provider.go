package provider

import (
	"errors"

	"github.com/odpf/stencil/core/schema"
	"github.com/odpf/stencil/formats/avro"
	"github.com/odpf/stencil/formats/json"
	"github.com/odpf/stencil/formats/protobuf"
)

type parseFn func([]byte) (schema.ParsedSchema, error)

type SchemaProvider struct {
	mapper map[string]parseFn
}

func (s *SchemaProvider) Parse(format string, data []byte) (schema.ParsedSchema, error) {
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
