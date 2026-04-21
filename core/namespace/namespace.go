// Package namespace provides types and operations for managing schema namespaces.
package namespace

import (
	"context"
	"time"
)

// Namespace represents a grouping of schemas with shared configuration.
type Namespace struct {
	ID            string
	Format        string
	Compatibility string
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Repository defines the persistence interface for namespaces.
type Repository interface {
	Create(context.Context, Namespace) (Namespace, error)
	Update(context.Context, Namespace) (Namespace, error)
	List(context.Context) ([]Namespace, error)
	Get(context.Context, string) (Namespace, error)
	Delete(context.Context, string) error
}
