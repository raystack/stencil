package store

import (
	"github.com/odpf/stencil/domain"
)

// Store is the interface that all database objects must implement.
type Store interface {
	domain.NamespaceRepository
	domain.SchemaRepository
}
