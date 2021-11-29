package json

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/stencil/server/domain"
	"github.com/odpf/stencil/server/schema"
	"go.uber.org/multierr"
)

const jsonFormat = "FORMAT_JSON"

type Schema struct {
	data []byte
}

func (s *Schema) Format() string {
	return jsonFormat
}

func (s *Schema) GetCanonicalValue() *domain.SchemaFile {
	id := uuid.NewSHA1(uuid.NameSpaceOID, s.data)
	return &domain.SchemaFile{
		ID:   id.String(),
		Data: s.data,
	}
}

func (s *Schema) verify(against schema.ParsedSchema) (*Schema, error) {
	prev, ok := against.(*Schema)
	if s.Format() == against.Format() && ok {
		return prev, nil
	}
	return nil, &runtime.HTTPStatusError{HTTPStatus: 400, Err: fmt.Errorf("current and prev schema formats(%s, %s) are different", s.Format(), against.Format())}
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
