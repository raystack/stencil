package storage

import (
	"github.com/odpf/stencil/domain"
	"github.com/odpf/stencil/server/namespace"
)

// Store is the interface that all database objects must implement.
type Store interface {
	namespace.Repository
	domain.SchemaRepository
}
