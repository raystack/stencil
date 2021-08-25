package api

import (
	"context"

	"github.com/odpf/stencil/server/api/v1/genproto"
	"github.com/odpf/stencil/server/snapshot"
)

// List returns list of snapshots. If filters applied it will return filtered snapshot list
func (a *API) List(ctx context.Context, req *genproto.ListSnapshotRequest) (*genproto.SnapshotList, error) {
	res := &genproto.SnapshotList{}
	list, err := a.Metadata.List(ctx, &snapshot.Snapshot{Namespace: req.Namespace, Name: req.Name, Version: req.Version, Latest: req.Latest})
	if err != nil {
		return res, err
	}
	for _, j := range list {
		res.Snapshots = append(res.Snapshots, fromSnapshotToProto(j))
	}
	return res, nil
}

// UpdateLatest marks specified snapshot as latest
func (a *API) UpdateLatest(ctx context.Context, req *genproto.UpdateLatestRequest) (*genproto.Snapshot, error) {
	var res *genproto.Snapshot
	st, err := a.Metadata.GetSnapshotByID(ctx, req.Id)
	if err != nil {
		return res, err
	}
	err = a.Metadata.UpdateLatestVersion(ctx, st)
	if err != nil {
		return res, err
	}
	return fromSnapshotToProto(st), nil
}
