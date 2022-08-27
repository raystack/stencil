package store

import (
	"github.com/odpf/stencil/core/namespace"
	"github.com/odpf/stencil/domain"
)

// Store is the interface that all database objects must implement.
type Store interface {
	namespace.NamespaceRepository
	domain.SchemaRepository
}
