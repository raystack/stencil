package storage

import (
	"github.com/odpf/stencil/server/namespace"
	"github.com/odpf/stencil/server/schema"
)

// Store is the interface that all database objects must implement.
type Store interface {
	namespace.Repository
	schema.Repository
}
