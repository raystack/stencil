package json

import (
	js "encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/stencil/server/schema"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func GetParsedSchema(data []byte) (schema.ParsedSchema, error) {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft2020
	sc, _ := compiler.Compile("https://json-schema.org/draft/2020-12/schema")
	var val interface{}
	if err := js.Unmarshal(data, &val); err != nil {
		return nil, &runtime.HTTPStatusError{HTTPStatus: http.StatusBadRequest, Err: err}
	}
	if err := sc.Validate(val); err != nil {
		return nil, &runtime.HTTPStatusError{HTTPStatus: http.StatusBadRequest, Err: err}
	}
	return &Schema{data: data}, nil
}
