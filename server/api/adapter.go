package api

import (
	"github.com/odpf/stencil/models"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
)

func fromProtoToSnapshot(g *stencilv1.Snapshot) *models.Snapshot {
	return &models.Snapshot{
		ID:        g.GetId(),
		Namespace: g.GetNamespace(),
		Name:      g.GetName(),
		Version:   g.GetVersion(),
		Latest:    g.GetLatest(),
	}
}

func fromSnapshotToProto(g *models.Snapshot) *stencilv1.Snapshot {
	return &stencilv1.Snapshot{
		Id:        g.ID,
		Namespace: g.Namespace,
		Name:      g.Name,
		Version:   g.Version,
		Latest:    g.Latest,
	}
}

func toRulesList(r *stencilv1.Checks) []string {
	var rules []string
	if r == nil {
		return rules
	}
	for _, rule := range r.Except {
		rules = append(rules, rule.String())
	}
	return rules
}

func toFileDownloadRequest(g *stencilv1.DownloadDescriptorRequest) *models.FileDownloadRequest {
	return &models.FileDownloadRequest{
		Namespace: g.Namespace,
		Name:      g.Name,
		Version:   g.Version,
		FullNames: g.GetFullnames(),
	}
}
