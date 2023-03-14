package json

import (
	"errors"

	"github.com/google/uuid"
	"github.com/goto/stencil/core/schema"
	"go.uber.org/multierr"
)

const jsonFormat = "FORMAT_JSON"

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
	return errors.New("Not implemented")
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
