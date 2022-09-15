package namespace

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}

func (s Service) Create(ctx context.Context, ns Namespace) (Namespace, error) {
	return s.repo.Create(ctx, ns)
}

func (s Service) Update(ctx context.Context, ns Namespace) (Namespace, error) {
	return s.repo.Update(ctx, ns)
}

func (s Service) List(ctx context.Context) ([]Namespace, error) {
	return s.repo.List(ctx)
}

func (s Service) Get(ctx context.Context, name string) (Namespace, error) {
	return s.repo.Get(ctx, name)
}

func (s Service) Delete(ctx context.Context, name string) error {
	return s.repo.Delete(ctx, name)
}
