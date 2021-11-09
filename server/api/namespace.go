package api

import (
	"context"

	"github.com/odpf/stencil/server/namespace"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createNamespaceRequestToNamespace(r *stencilv1.CreateNamespaceRequest) namespace.Namespace {
	return namespace.Namespace{
		ID:          r.GetId(),
		Format:      r.GetFormat().String(),
		Description: r.GetDescription(),
	}
}

func namespaceToProto(ns namespace.Namespace) *stencilv1.Namespace {
	return &stencilv1.Namespace{
		Id:          ns.ID,
		Format:      stencilv1.Schema_Format(stencilv1.Schema_Format_value[ns.Format]),
		Description: ns.Description,
		CreatedAt:   timestamppb.New(ns.CreatedAt),
		UpdatedAt:   timestamppb.New(ns.UpdatedAt),
	}
}

// CreateNamespace handler for creating namespace
func (a *API) CreateNamespace(ctx context.Context, in *stencilv1.CreateNamespaceRequest) (*stencilv1.CreateNamespaceResponse, error) {
	ns := createNamespaceRequestToNamespace(in)
	newNamespace, err := a.Namespace.Create(ctx, ns)
	return &stencilv1.CreateNamespaceResponse{Namespace: namespaceToProto(newNamespace)}, err
}

func (a *API) UpdateNamespace(ctx context.Context, in *stencilv1.UpdateNamespaceRequest) (*stencilv1.UpdateNamespaceResponse, error) {
	ns, err := a.Namespace.Update(ctx, namespace.Namespace{ID: in.GetId(), Format: in.GetFormat().String(), Description: in.GetDescription()})
	return &stencilv1.UpdateNamespaceResponse{Namespace: namespaceToProto(ns)}, err
}

func (a *API) GetNamespace(ctx context.Context, in *stencilv1.GetNamespaceRequest) (*stencilv1.GetNamespaceResponse, error) {
	namespace, err := a.Namespace.Get(ctx, in.GetId())
	return &stencilv1.GetNamespaceResponse{Namespace: namespaceToProto(namespace)}, err
}

// ListNamespaces handler for returning list of available namespaces
func (a *API) ListNamespaces(ctx context.Context, in *stencilv1.ListNamespacesRequest) (*stencilv1.ListNamespacesResponse, error) {
	namespaces, err := a.Namespace.List(ctx)
	return &stencilv1.ListNamespacesResponse{Namespaces: namespaces}, err
}

func (a *API) DeleteNamespace(ctx context.Context, in *stencilv1.DeleteNamespaceRequest) (*stencilv1.DeleteNamespaceResponse, error) {
	err := a.Namespace.Delete(ctx, in.GetId())
	message := "success"
	if err != nil {
		message = "failed"
	}
	return &stencilv1.DeleteNamespaceResponse{Message: message}, err
}
