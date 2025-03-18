package postgres

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/raystack/stencil/core/namespace"
)

const namespaceListQuery = `
SELECT id, format, compatibility from namespaces
`

const namespaceGetQuery = `
SELECT * from namespaces where id=$1
`

const namespaceDeleteQuery = `
DELETE from namespaces where id=$1
`

const namespaceUpdateQuery = `
UPDATE namespaces SET format=$2,compatibility=$3,description=$4,updated_at=now()
WHERE id = $1
RETURNING *
`

const namespaceInsertQuery = `
INSERT INTO namespaces (id, format, compatibility, description, created_at, updated_at)
    VALUES ($1, $2, $3, $4, now(), now())
RETURNING *
`

type NamespaceRepository struct {
	db *DB
}

func NewNamespaceRepository(dbc *DB) *NamespaceRepository {
	return &NamespaceRepository{
		db: dbc,
	}
}

func (r *NamespaceRepository) Create(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error) {
	newNamespace := namespace.Namespace{}
	err := pgxscan.Get(ctx, r.db, &newNamespace, namespaceInsertQuery, ns.ID, ns.Format, ns.Compatibility, ns.Description)
	return newNamespace, wrapError(err, "%s", ns.ID)
}

func (r *NamespaceRepository) Update(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error) {
	newNamespace := namespace.Namespace{}
	err := pgxscan.Get(ctx, r.db, &newNamespace, namespaceUpdateQuery, ns.ID, ns.Format, ns.Compatibility, ns.Description)
	return newNamespace, wrapError(err, "%s", ns.ID)
}

func (r *NamespaceRepository) Get(ctx context.Context, id string) (namespace.Namespace, error) {
	newNamespace := namespace.Namespace{}
	err := pgxscan.Get(ctx, r.db, &newNamespace, namespaceGetQuery, id)
	return newNamespace, wrapError(err, "%s", id)
}

func (r *NamespaceRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, namespaceDeleteQuery, id)
	r.db.Exec(ctx, deleteOrphanedData)
	return wrapError(err, "%s", id)
}

func (r *NamespaceRepository) List(ctx context.Context) ([]namespace.Namespace, error) {
	var namespaces []namespace.Namespace
	err := pgxscan.Select(ctx, r.db, &namespaces, namespaceListQuery)
	return namespaces, wrapError(err, "")
}
