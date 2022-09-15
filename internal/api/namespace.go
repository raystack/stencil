package api

import (
	"context"

	"github.com/odpf/stencil/core/namespace"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createNamespaceRequestToNamespace(r *stencilv1beta1.CreateNamespaceRequest) namespace.Namespace {
	return namespace.Namespace{
		ID:            r.GetId(),
		Format:        r.GetFormat().String(),
		Compatibility: r.GetCompatibility().String(),
		Description:   r.GetDescription(),
	}
}

func namespaceToProto(ns namespace.Namespace) *stencilv1beta1.Namespace {
	return &stencilv1beta1.Namespace{
		Id:            ns.ID,
		Format:        stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[ns.Format]),
		Compatibility: stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[ns.Compatibility]),
		Description:   ns.Description,
		CreatedAt:     timestamppb.New(ns.CreatedAt),
		UpdatedAt:     timestamppb.New(ns.UpdatedAt),
	}
}

// CreateNamespace handler for creating namespace
func (a *API) CreateNamespace(ctx context.Context, in *stencilv1beta1.CreateNamespaceRequest) (*stencilv1beta1.CreateNamespaceResponse, error) {
	ns := createNamespaceRequestToNamespace(in)
	newNamespace, err := a.namespace.Create(ctx, ns)
	return &stencilv1beta1.CreateNamespaceResponse{Namespace: namespaceToProto(newNamespace)}, err
}

func (a *API) UpdateNamespace(ctx context.Context, in *stencilv1beta1.UpdateNamespaceRequest) (*stencilv1beta1.UpdateNamespaceResponse, error) {
	ns, err := a.namespace.Update(ctx, namespace.Namespace{ID: in.GetId(), Format: in.GetFormat().String(), Compatibility: in.GetCompatibility().String(), Description: in.GetDescription()})
	return &stencilv1beta1.UpdateNamespaceResponse{Namespace: namespaceToProto(ns)}, err
}

func (a *API) GetNamespace(ctx context.Context, in *stencilv1beta1.GetNamespaceRequest) (*stencilv1beta1.GetNamespaceResponse, error) {
	namespace, err := a.namespace.Get(ctx, in.GetId())
	return &stencilv1beta1.GetNamespaceResponse{Namespace: namespaceToProto(namespace)}, err
}

// ListNamespaces handler for returning list of available namespaces
func (a *API) ListNamespaces(ctx context.Context, in *stencilv1beta1.ListNamespacesRequest) (*stencilv1beta1.ListNamespacesResponse, error) {
	namespaces, err := a.namespace.List(ctx)
	var nsp []*stencilv1beta1.Namespace
	for _, n := range namespaces {
		nsp = append(nsp, namespaceToProto(n))
	}
	return &stencilv1beta1.ListNamespacesResponse{Namespaces: nsp}, err
}

func (a *API) DeleteNamespace(ctx context.Context, in *stencilv1beta1.DeleteNamespaceRequest) (*stencilv1beta1.DeleteNamespaceResponse, error) {
	err := a.namespace.Delete(ctx, in.GetId())
	message := "success"
	if err != nil {
		message = "failed"
	}
	return &stencilv1beta1.DeleteNamespaceResponse{Message: message}, err
}
