package api

import (
	"github.com/odpf/stencil/server/api/v1/pb"
	"github.com/odpf/stencil/server/snapshot"
)

func fromProtoToSnapshot(g *pb.Snapshot) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        g.GetId(),
		Namespace: g.GetNamespace(),
		Name:      g.GetName(),
		Version:   g.GetVersion(),
		Latest:    g.GetLatest(),
	}
}

func fromSnapshotToProto(g *snapshot.Snapshot) *pb.Snapshot {
	return &pb.Snapshot{
		Id:        g.ID,
		Namespace: g.Namespace,
		Name:      g.Name,
		Version:   g.Version,
		Latest:    g.Latest,
	}
}
