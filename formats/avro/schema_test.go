package avro_test

import (
	"testing"

	"github.com/goto/stencil/formats/avro"
	"github.com/stretchr/testify/assert"
)

// adding basic tests only since these compatibility checks already tested in avro lib
func TestAvroCompatibility(t *testing.T) {
	for _, test := range []struct {
		name         string
		writerSchema string
		readerSchema string
		isError      bool
	}{
		{"field addition with default allowed",
			`{
				"type": "record",
				"name": "myrecord",
				"fields": [{ "type": "string", "name": "f1" }]
			}`,

			`{
				"type": "record",
				"name": "myrecord",
				"fields": [
					{ "type": "string", "name": "f1" },
					{ "type": "string", "name": "f2", "default": "some" }
				]
			}`,
			false},
		{"field addition without default is not allowed",
			`{
				"type": "record",
				"name": "myrecord",
				"fields": [{ "type": "string", "name": "f1" }]
			}`,

			`{
				"type": "record",
				"name": "myrecord",
				"fields": [
					{ "type": "string", "name": "f1" },
					{ "type": "string", "name": "f2" }
				]
			}`,
			true},
		{"making field type as union allowed",
			`{
				"type": "record",
				"name": "myrecord",
				"fields": [{ "type": "string", "name": "f1" }]
			}`,

			`{
				"type": "record",
				"name": "myrecord",
				"fields": [
					{ "type": [null, "string"], "name": "f1" }
				]
			}`,
			false},
	} {
		t.Run(test.name, func(t *testing.T) {
			prev, err := avro.ParseSchema([]byte(test.writerSchema))
			assert.NoError(t, err)
			current, err := avro.ParseSchema([]byte(test.readerSchema))
			assert.NoError(t, err)
			err = current.IsBackwardCompatible(prev)
			if test.isError {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
			err = prev.IsForwardCompatible(current)
			if test.isError {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseSchema(t *testing.T) {
	t.Run("should return error if avro schema is not valid", func(t *testing.T) {
		s := `{
			"type": "invalid",
			"name": "test",
			"fields": "invalid"
		}`
		_, err := avro.ParseSchema([]byte(s))
		assert.NotNil(t, err)
	})
}
