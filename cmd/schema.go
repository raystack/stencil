package cmd

import (
	"context"

	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
)

func SchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "schema",
		Aliases: []string{"schemas"},
		Short:   "Manage schemas",
		Long:    "Work with schemas.",
		Annotations: map[string]string{
			"group": "core",
		},
	}

	cmd.AddCommand(createSchemaCmd())
	cmd.AddCommand(listSchemaCmd())
	cmd.AddCommand(infoSchemaCmd())
	cmd.AddCommand(versionSchemaCmd())
	cmd.AddCommand(printSchemaCmd())
	cmd.AddCommand(downloadSchemaCmd())
	cmd.AddCommand(checkSchemaCmd())
	cmd.AddCommand(editSchemaCmd())
	cmd.AddCommand(deleteSchemaCmd())
	cmd.AddCommand(diffSchemaCmd())
	cmd.AddCommand(graphSchemaCmd())

	return cmd
}

func fetchSchemaAndMeta(client stencilv1beta1.StencilServiceClient, version int32, namespaceID, schemaID string) ([]byte, *stencilv1beta1.GetSchemaMetadataResponse, error) {
	var req stencilv1beta1.GetSchemaRequest
	var reqLatest stencilv1beta1.GetLatestSchemaRequest
	var data []byte

	ctx := context.Background()

	if version != 0 {
		req.NamespaceId = namespaceID
		req.SchemaId = schemaID
		req.VersionId = version
		res, err := client.GetSchema(ctx, &req)
		if err != nil {
			return nil, nil, err
		}
		data = res.GetData()
	} else {
		reqLatest.NamespaceId = namespaceID
		reqLatest.SchemaId = schemaID
		res, err := client.GetLatestSchema(ctx, &reqLatest)
		if err != nil {
			return nil, nil, err
		}
		data = res.GetData()
	}

	reqMeta := stencilv1beta1.GetSchemaMetadataRequest{
		NamespaceId: namespaceID,
		SchemaId:    schemaID,
	}
	meta, err := client.GetSchemaMetadata(context.Background(), &reqMeta)

	if err != nil {
		return nil, nil, err
	}

	return data, meta, nil
}

func fetchMeta(client stencilv1beta1.StencilServiceClient, namespace string, schema string) (*stencilv1beta1.GetSchemaMetadataResponse, error) {
	req := stencilv1beta1.GetSchemaMetadataRequest{
		NamespaceId: namespace,
		SchemaId:    schema,
	}
	return client.GetSchemaMetadata(context.Background(), &req)
}
