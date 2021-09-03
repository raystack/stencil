package api

import (
	"context"

	"github.com/odpf/stencil/server/api/v1/pb"
	"github.com/odpf/stencil/server/snapshot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListSnapshots returns list of snapshots. If filters applied it will return filtered snapshot list
func (a *API) ListSnapshots(ctx context.Context, req *pb.ListSnapshotsRequest) (*pb.ListSnapshotsResponse, error) {
	res := &pb.ListSnapshotsResponse{}
	list, err := a.Metadata.List(ctx, &snapshot.Snapshot{Namespace: req.Namespace, Name: req.Name, Version: req.Version, Latest: req.Latest})
	if err != nil {
		return res, err
	}
	for _, j := range list {
		res.Snapshots = append(res.Snapshots, fromSnapshotToProto(j))
	}
	return res, nil
}

// PromoteSnapshot marks specified snapshot as latest
func (a *API) PromoteSnapshot(ctx context.Context, req *pb.PromoteSnapshotRequest) (*pb.Snapshot, error) {
	var res *pb.Snapshot
	st, err := a.Metadata.GetSnapshotByID(ctx, req.Id)
	if err != nil {
		if err == snapshot.ErrNotFound {
			return res, status.Error(codes.NotFound, err.Error())
		}
		return res, status.Error(codes.Internal, err.Error())
	}
	err = a.Metadata.UpdateLatestVersion(ctx, st)
	if err != nil {
		return res, err
	}
	return fromSnapshotToProto(st), nil
}
