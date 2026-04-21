package namespace

import (
	"context"
)

// Service provides namespace management operations.
type Service struct {
	repo Repository
}

// NewService creates a new namespace service with the given repository.
func NewService(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}

// Create stores a new namespace.
func (s Service) Create(ctx context.Context, ns Namespace) (Namespace, error) {
	return s.repo.Create(ctx, ns)
}

// Update modifies an existing namespace.
func (s Service) Update(ctx context.Context, ns Namespace) (Namespace, error) {
	return s.repo.Update(ctx, ns)
}

// List returns all namespaces.
func (s Service) List(ctx context.Context) ([]Namespace, error) {
	return s.repo.List(ctx)
}

// Get retrieves a namespace by name.
func (s Service) Get(ctx context.Context, name string) (Namespace, error) {
	return s.repo.Get(ctx, name)
}

// Delete removes a namespace by name.
func (s Service) Delete(ctx context.Context, name string) error {
	return s.repo.Delete(ctx, name)
}
