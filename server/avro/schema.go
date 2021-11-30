package avro

import (
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	av "github.com/hamba/avro"
	"github.com/odpf/stencil/server/domain"
	"github.com/odpf/stencil/server/schema"
	"go.uber.org/multierr"
)

const avroFormat = "FORMAT_AVRO"

type Schema struct {
	data []byte
	sc   av.Schema
}

func (s *Schema) Format() string {
	return avroFormat
}

func (s *Schema) GetCanonicalValue() *domain.SchemaFile {
	id := s.sc.Fingerprint()
	return &domain.SchemaFile{
		ID:   string(id[:]),
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
	prev, err := s.verify(against)
	if err != nil {
		return err
	}
	c := av.NewSchemaCompatibility()
	return c.Compatible(s.sc, prev.sc)
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
