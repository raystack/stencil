package avro

import (
	"net/http"

	"github.com/goto/stencil/core/schema"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	av "github.com/hamba/avro"
)

// ParseSchema parses avro schema bytes into ParsedSchema
func ParseSchema(data []byte) (schema.ParsedSchema, error) {
	sc, err := av.Parse(string(data))
	if err != nil {
		return nil, &runtime.HTTPStatusError{HTTPStatus: http.StatusBadRequest, Err: err}
	}
	return &Schema{sc: sc, data: data}, nil
}
