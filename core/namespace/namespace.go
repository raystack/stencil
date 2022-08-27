package namespace

import (
	"context"
	"time"
)

// Namespace model
type Namespace struct {
	ID            string
	Format        string
	Compatibility string
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NamespaceService namespace CRUD service
type NamespaceService interface {
	Create(ctx context.Context, ns Namespace) (Namespace, error)
	Update(ctx context.Context, ns Namespace) (Namespace, error)
	List(ctx context.Context) ([]string, error)
	Get(ctx context.Context, name string) (Namespace, error)
	Delete(ctx context.Context, name string) error
}

// NamespaceRepository for namespace
type NamespaceRepository interface {
	CreateNamespace(context.Context, Namespace) (Namespace, error)
	UpdateNamespace(context.Context, Namespace) (Namespace, error)
	ListNamespaces(context.Context) ([]string, error)
	GetNamespace(context.Context, string) (Namespace, error)
	DeleteNamespace(context.Context, string) error
}
