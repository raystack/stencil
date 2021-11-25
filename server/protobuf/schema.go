package protobuf

import (
	"errors"

	"github.com/google/uuid"
	"github.com/odpf/stencil/server/schema"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const protobufFormat = "FORMAT_PROTOBUF"

type Schema struct {
	*protoregistry.Files
	isValid bool
	data    []byte
}

func (s *Schema) Validate() error {
	if s.Files != nil {
		return nil
	}
	return errors.New("not valid schema")
}

func (s *Schema) Format() string {
	return protobufFormat
}

func (s *Schema) GetCanonicalValue() *schema.SchemaFile {
	id := uuid.NewSHA1(uuid.NameSpaceOID, s.data)
	return &schema.SchemaFile{
		ID:     id.String(),
		Types:  getAllMessages(s.Files),
		Data:   s.data,
		Fields: getAllFields(s.Files),
	}
}

// IsBackwardCompatible checks backward compatibility against given schema
// Allowed changes: field addition
// Disallowed changes: field type change, tag number change, label change
func (s *Schema) IsBackwardCompatible(against schema.ParsedSchema) error {
	prev, ok := against.(*Schema)
	if against.Format() != protobufFormat && !ok {
		return errors.New("different schema formats")
	}
	return compareSchemas(s.Files, prev.Files, backwardCompatibility)
}

// IsForwardCompatible for protobuf forward compatible is same as backward compatible
func (s *Schema) IsForwardCompatible(against schema.ParsedSchema) error {
	prev, ok := against.(*Schema)
	if against.Format() != protobufFormat && !ok {
		return errors.New("different schema formats")
	}
	return compareSchemas(s.Files, prev.Files, forwardCompatibility)
}

// IsFullCompatible for protobuf forward compatible is same as backward compatible
func (s *Schema) IsFullCompatible(against schema.ParsedSchema) error {
	prev, ok := against.(*Schema)
	if against.Format() != protobufFormat && !ok {
		return errors.New("different schema formats")
	}
	return compareSchemas(s.Files, prev.Files, forwardCompatibility)
}
