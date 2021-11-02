package schema

import (
	"context"
	"time"

	"github.com/odpf/stencil/server/namespace"
)

// Schema model
type Schema struct {
	ID            string
	Authority     string
	Format        string
	Compatibility string
	Description   string
	NamespaceID   string
	Data          []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Repository for Schema
type Repository interface {
	CreateSchema(context.Context, *Schema) (*Schema, error)
	// UpdateSchema(context.Context, Schema) (Schema, error)
	ListSchemas(context.Context, string) ([]string, error)
	// GetSchema(context.Context, string) (Schema, error)
	// DeleteSchema(context.Context, string) error
}

type Service struct {
	Repo         Repository
	NamespaceSvc namespace.Service
}

func (s *Service) isNamespaceExist(ctx context.Context, id string) error {
	_, err := s.NamespaceSvc.GetNamespace(ctx, id)
	return err
}

func (s *Service) CreateSchema(ctx context.Context, sc *Schema) (*Schema, error) {
	err := s.isNamespaceExist(ctx, sc.NamespaceID)
	if err != nil {
		return nil, err
	}
	return s.Repo.CreateSchema(ctx, sc)
}

func (s *Service) ListSchemas(ctx context.Context, namespaceID string) ([]string, error) {
	err := s.isNamespaceExist(ctx, namespaceID)
	if err != nil {
		return nil, err
	}
	return s.Repo.ListSchemas(ctx, namespaceID)
}
