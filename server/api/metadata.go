package api

import (
	"context"

	"github.com/odpf/stencil/models"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListSnapshots returns list of snapshots. If filters applied it will return filtered snapshot list
func (a *API) ListSnapshots(ctx context.Context, req *stencilv1.ListSnapshotsRequest) (*stencilv1.ListSnapshotsResponse, error) {
	res := &stencilv1.ListSnapshotsResponse{}
	list, err := a.Metadata.List(ctx, &models.Snapshot{Namespace: req.Namespace, Name: req.Name, Version: req.Version, Latest: req.Latest})
	if err != nil {
		return res, err
	}
	for _, j := range list {
		res.Snapshots = append(res.Snapshots, fromSnapshotToProto(j))
	}
	return res, nil
}

// PromoteSnapshot marks specified snapshot as latest
func (a *API) PromoteSnapshot(ctx context.Context, req *stencilv1.PromoteSnapshotRequest) (*stencilv1.PromoteSnapshotResponse, error) {
	st, err := a.Metadata.GetSnapshotByID(ctx, req.Id)
	if err != nil {
		if err == models.ErrSnapshotNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = a.Metadata.UpdateLatestVersion(ctx, st)
	if err != nil {
		return nil, err
	}
	return &stencilv1.PromoteSnapshotResponse{
		Snapshot: fromSnapshotToProto(st),
	}, nil
}
