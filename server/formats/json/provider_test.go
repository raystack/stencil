package json_test

import (
	"testing"

	"github.com/odpf/stencil/server/formats/json"
	"github.com/stretchr/testify/assert"
)

func TestSchemaValidation(t *testing.T) {
	for _, test := range []struct {
		name    string
		schema  string
		isError bool
	}{
		{"should return error if schema is not json", `{invalid_json,}`, true},
		{"should return error if schema is not valid json", `{
			"$id": "https://example.com/address.schema.json",
			"type": "object",
			"properties": {
				"f1": {
					"type": "string"
				}
			},
			"required": "this is supposed to array"
		}`, true},
		{"should return nil if schema is valid json", `{
			"$id": "https://example.com/address.schema.json",
			"type": "object",
			"properties": {
				"f1": {
					"type": "string"
				}
			},
			"required": ["f1"]
		}`, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := json.GetParsedSchema([]byte(test.schema))
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
