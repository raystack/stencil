package json

import (
	js "encoding/json"

	"connectrpc.com/connect"
	"github.com/raystack/stencil/core/schema"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func GetParsedSchema(data []byte) (schema.ParsedSchema, error) {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft2020
	sc, _ := compiler.Compile("https://json-schema.org/draft/2020-12/schema")
	var val interface{}
	if err := js.Unmarshal(data, &val); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if err := sc.Validate(val); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return &Schema{data: data}, nil
}
