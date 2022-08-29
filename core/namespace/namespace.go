package namespace

import (
	"context"
	"time"
)

type Namespace struct {
	ID            string
	Format        string
	Compatibility string
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Repository interface {
	Create(context.Context, Namespace) (Namespace, error)
	Update(context.Context, Namespace) (Namespace, error)
	List(context.Context) ([]string, error)
	Get(context.Context, string) (Namespace, error)
	Delete(context.Context, string) error
}
