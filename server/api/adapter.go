package api

import (
	"github.com/odpf/stencil/server/api/v1/genproto"
	"github.com/odpf/stencil/server/snapshot"
)

func fromProtoToSnapshot(g *genproto.Snapshot) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        g.GetId(),
		Namespace: g.GetNamespace(),
		Name:      g.GetName(),
		Version:   g.GetVersion(),
		Latest:    g.GetLatest(),
	}
}

func fromSnapshotToProto(g *snapshot.Snapshot) *genproto.Snapshot {
	return &genproto.Snapshot{
		Id:        g.ID,
		Namespace: g.Namespace,
		Name:      g.Name,
		Version:   g.Version,
		Latest:    g.Latest,
	}
}
