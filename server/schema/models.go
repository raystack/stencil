package schema

import (
	"github.com/odpf/stencil/server/domain"
)

type ParsedSchema interface {
	IsBackwardCompatible(ParsedSchema) error
	IsForwardCompatible(ParsedSchema) error
	IsFullCompatible(ParsedSchema) error
	Format() string
	GetCanonicalValue() *domain.SchemaFile
}

type SchemaProvider interface {
	ParseSchema(format string, data []byte) (ParsedSchema, error)
}
