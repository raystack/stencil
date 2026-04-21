package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/raystack/stencil/core/namespace"
	stencilv1beta1 "github.com/raystack/stencil/gen/raystack/stencil/v1beta1"
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
func (a *API) CreateNamespace(ctx context.Context, req *connect.Request[stencilv1beta1.CreateNamespaceRequest]) (*connect.Response[stencilv1beta1.CreateNamespaceResponse], error) {
	ns := createNamespaceRequestToNamespace(req.Msg)
	newNamespace, err := a.namespace.Create(ctx, ns)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&stencilv1beta1.CreateNamespaceResponse{Namespace: namespaceToProto(newNamespace)}), nil
}

func (a *API) UpdateNamespace(ctx context.Context, req *connect.Request[stencilv1beta1.UpdateNamespaceRequest]) (*connect.Response[stencilv1beta1.UpdateNamespaceResponse], error) {
	ns, err := a.namespace.Update(ctx, namespace.Namespace{ID: req.Msg.GetId(), Format: req.Msg.GetFormat().String(), Compatibility: req.Msg.GetCompatibility().String(), Description: req.Msg.GetDescription()})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&stencilv1beta1.UpdateNamespaceResponse{Namespace: namespaceToProto(ns)}), nil
}

func (a *API) GetNamespace(ctx context.Context, req *connect.Request[stencilv1beta1.GetNamespaceRequest]) (*connect.Response[stencilv1beta1.GetNamespaceResponse], error) {
	namespace, err := a.namespace.Get(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&stencilv1beta1.GetNamespaceResponse{Namespace: namespaceToProto(namespace)}), nil
}

// ListNamespaces handler for returning list of available namespaces
func (a *API) ListNamespaces(ctx context.Context, req *connect.Request[stencilv1beta1.ListNamespacesRequest]) (*connect.Response[stencilv1beta1.ListNamespacesResponse], error) {
	namespaces, err := a.namespace.List(ctx)
	if err != nil {
		return nil, err
	}
	var nsp []*stencilv1beta1.Namespace
	for _, n := range namespaces {
		nsp = append(nsp, namespaceToProto(n))
	}
	return connect.NewResponse(&stencilv1beta1.ListNamespacesResponse{Namespaces: nsp}), nil
}

func (a *API) DeleteNamespace(ctx context.Context, req *connect.Request[stencilv1beta1.DeleteNamespaceRequest]) (*connect.Response[stencilv1beta1.DeleteNamespaceResponse], error) {
	err := a.namespace.Delete(ctx, req.Msg.GetId())
	message := "success"
	if err != nil {
		message = "failed"
	}
	return connect.NewResponse(&stencilv1beta1.DeleteNamespaceResponse{Message: message}), err
}
