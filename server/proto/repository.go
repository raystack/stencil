package proto

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/odpf/stencil/server/snapshot"
	"github.com/odpf/stencil/server/store"
)

const snapshotInsertQuery = `
WITH ss(id) as (
	INSERT INTO snapshots (namespace, name, version)
	VALUES ($1, $2, $3) ON CONFLICT DO NOTHING
	RETURNING snapshots.id
)
SELECT COALESCE(
		(
			select ss.id
			from ss
		),
		(
			select id
			from snapshots
			where namespace = $1
				and name = $2
				and version = $3
		)
	)`

const fileInsertQuery = `
WITH file_insert(id) as (
	INSERT INTO protobuf_files (search_data, data)
	VALUES ($2, $3) ON CONFLICT DO NOTHING
	RETURNING id
),
file(id) as (
	SELECT COALESCE(
			(
				SELECT id
				FROM file_insert
			),
			(
				select id
				from protobuf_files
				where search_data = $2
					and data = $3
			)
		)
)
INSERT INTO snapshots_protobuf_files(snapshot_id, file_id)
SELECT $1,file.id from file`

const getDataForSpecificMessages = `
WITH files as (
	SELECT distinct(jsonb_array_elements_text(pf.search_data->'dependencies'))
	from protobuf_files as pf
		join snapshots_protobuf_files as spf on pf.id = spf.file_id
		join snapshots s on s.id = spf.snapshot_id
	WHERE spf.snapshot_id = $1 AND pf.search_data->'messages' ?| $2
)
	SELECT pf.data
	from protobuf_files as pf
		join snapshots_protobuf_files as spf on pf.id = spf.file_id
		join snapshots s on s.id = spf.snapshot_id
	WHERE spf.snapshot_id = $1 and pf.search_data->>'path' in (select * from files);
`

const getWholeFDS = `
SELECT pf.data
	from protobuf_files as pf
		join snapshots_protobuf_files as spf on pf.id = spf.file_id
		join snapshots s on s.id = spf.snapshot_id
	WHERE spf.snapshot_id = $1
`

// Repository DB access layer
type Repository struct {
	db *store.DB
}

// Put inserts fileDescriptorset information in DB
func (r *Repository) Put(ctx context.Context, snapshot *snapshot.Snapshot, dbFiles []*ProtobufDBFile) error {
	return r.db.Pool.BeginFunc(ctx, func(t pgx.Tx) error {
		err := t.QueryRow(ctx, snapshotInsertQuery, snapshot.Namespace, snapshot.Name, snapshot.Version).Scan(&snapshot.ID)
		if err != nil {
			return err
		}

		batch := &pgx.Batch{}
		for _, file := range dbFiles {
			batch.Queue(fileInsertQuery, snapshot.ID, file.SearchData, file.Data)
		}
		res := t.SendBatch(ctx, batch)
		for i := 0; i < len(dbFiles); i++ {
			_, err = res.Exec()
			if err != nil {
				return err
			}
		}
		err = res.Close()
		return err
	})
}

// Get Fullycontained descriptorset file given list of fully qualified message names.
// If message names are empty then whole fileDescriptorSet data returned
func (r *Repository) Get(ctx context.Context, snapshot *snapshot.Snapshot, names []string) ([][]byte, error) {
	var totalData [][]byte
	var err error
	if len(names) > 0 {
		err = pgxscan.Select(ctx, r.db, &totalData, getDataForSpecificMessages, snapshot.ID, names)
	} else {
		err = pgxscan.Select(ctx, r.db, &totalData, getWholeFDS, snapshot.ID)
	}
	return totalData, err
}

// NewProtoRepository create instance repo
func NewProtoRepository(db *store.DB) *Repository {
	return &Repository{
		db,
	}
}
