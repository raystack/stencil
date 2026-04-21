package avro

import (
	"connectrpc.com/connect"
	av "github.com/hamba/avro"
	"github.com/raystack/stencil/core/schema"
)

// ParseSchema parses avro schema bytes into ParsedSchema
func ParseSchema(data []byte) (schema.ParsedSchema, error) {
	sc, err := av.Parse(string(data))
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return &Schema{sc: sc, data: data}, nil
}
