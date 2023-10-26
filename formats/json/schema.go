package json

import (
	"github.com/google/uuid"
	"github.com/raystack/stencil/core/schema"
	"github.com/raystack/stencil/pkg/logger"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader" // imported to compile http references in json schema
	"go.uber.org/multierr"
)

const jsonFormat = "FORMAT_JSON"
const schemaURI = "sample_schema"

type Schema struct {
	data []byte
}

func (s *Schema) Format() string {
	return jsonFormat
}

func (s *Schema) GetCanonicalValue() *schema.SchemaFile {
	id := uuid.NewSHA1(uuid.NameSpaceOID, s.data)
	return &schema.SchemaFile{
		ID:   id.String(),
		Data: s.data,
	}
}

// IsBackwardCompatible checks backward compatibility against given schema
func (s *Schema) IsBackwardCompatible(against schema.ParsedSchema) error {
	sc, err := jsonschema.CompileString(schemaURI, string(s.data))
	if err != nil {
		logger.Logger.Warn("unable to compile schema to check for backward compatibility")
		return err
	}
	againstSchema, err := jsonschema.CompileString(schemaURI, string(against.GetCanonicalValue().Data))
	if err != nil {
		logger.Logger.Warn("unable to compile against schema to check for backward compatibility")
		return err
	}
	jsonSchemaMap := exploreSchema(sc)
	againstJsonSchemaMap := exploreSchema(againstSchema)
	return compareSchemas(againstJsonSchemaMap, jsonSchemaMap, backwardCompatibility,
		[]SchemaCompareCheck{CheckPropertyDeleted, TypeCheckExecutor(StandardTypeChecks)}, []SchemaCheck{CheckAdditionalProperties})
}

// IsForwardCompatible checks backward compatibility against given schema
func (s *Schema) IsForwardCompatible(against schema.ParsedSchema) error {
	return against.IsBackwardCompatible(s)
}

// IsFullCompatible checks for forward compatibility
func (s *Schema) IsFullCompatible(against schema.ParsedSchema) error {
	forwardErr := s.IsForwardCompatible(against)
	backwardErr := s.IsBackwardCompatible(against)
	return multierr.Combine(forwardErr, backwardErr)
}
