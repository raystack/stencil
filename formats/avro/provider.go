package avro

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	av "github.com/hamba/avro"
	"github.com/odpf/stencil/server/schema"
)

//ParseSchema parses avro schema bytes into ParsedSchema
func ParseSchema(data []byte) (schema.ParsedSchema, error) {
	sc, err := av.Parse(string(data))
	if err != nil {
		return nil, &runtime.HTTPStatusError{HTTPStatus: http.StatusBadRequest, Err: err}
	}
	return &Schema{sc: sc, data: data}, nil
}
