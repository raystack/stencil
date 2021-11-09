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

// Repository for namespace
type Repository interface {
	CreateNamespace(context.Context, Namespace) (Namespace, error)
	UpdateNamespace(context.Context, Namespace) (Namespace, error)
	ListNamespaces(context.Context) ([]string, error)
	GetNamespace(context.Context, string) (Namespace, error)
	DeleteNamespace(context.Context, string) error
}

type Service struct {
	Repo Repository
}

func (s Service) Create(ctx context.Context, ns Namespace) (Namespace, error) {
	return s.Repo.CreateNamespace(ctx, ns)
}

func (s Service) Update(ctx context.Context, ns Namespace) (Namespace, error) {
	return s.Repo.UpdateNamespace(ctx, ns)
}

func (s Service) List(ctx context.Context) ([]string, error) {
	return s.Repo.ListNamespaces(ctx)
}

func (s Service) Get(ctx context.Context, name string) (Namespace, error) {
	return s.Repo.GetNamespace(ctx, name)
}

func (s Service) Delete(ctx context.Context, name string) error {
	return s.Repo.DeleteNamespace(ctx, name)
}
