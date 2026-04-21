package cmd

import (
	"context"

	"connectrpc.com/connect"
	stencilv1beta1 "github.com/raystack/stencil/gen/raystack/stencil/v1beta1"
	stencilv1beta1connect "github.com/raystack/stencil/gen/raystack/stencil/v1beta1/stencilv1beta1connect"
	"github.com/spf13/cobra"
)

func SchemaCmd(cdk *CDK) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "schema",
		Aliases: []string{"schemas"},
		Short:   "Manage schemas",
		Long:    "Work with schemas.",
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(createSchemaCmd(cdk))
	cmd.AddCommand(listSchemaCmd(cdk))
	cmd.AddCommand(infoSchemaCmd(cdk))
	cmd.AddCommand(versionSchemaCmd(cdk))
	cmd.AddCommand(printSchemaCmd(cdk))
	cmd.AddCommand(downloadSchemaCmd(cdk))
	cmd.AddCommand(checkSchemaCmd(cdk))
	cmd.AddCommand(editSchemaCmd(cdk))
	cmd.AddCommand(deleteSchemaCmd(cdk))
	cmd.AddCommand(diffSchemaCmd(cdk))
	cmd.AddCommand(graphSchemaCmd(cdk))

	return cmd
}

func fetchSchemaAndMeta(client stencilv1beta1connect.StencilServiceClient, version int32, namespaceID, schemaID string) ([]byte, *stencilv1beta1.GetSchemaMetadataResponse, error) {
	var data []byte

	ctx := context.Background()

	if version != 0 {
		res, err := client.GetSchema(ctx, connect.NewRequest(&stencilv1beta1.GetSchemaRequest{
			NamespaceId: namespaceID,
			SchemaId:    schemaID,
			VersionId:   version,
		}))
		if err != nil {
			return nil, nil, err
		}
		data = res.Msg.GetData()
	} else {
		res, err := client.GetLatestSchema(ctx, connect.NewRequest(&stencilv1beta1.GetLatestSchemaRequest{
			NamespaceId: namespaceID,
			SchemaId:    schemaID,
		}))
		if err != nil {
			return nil, nil, err
		}
		data = res.Msg.GetData()
	}

	metaRes, err := client.GetSchemaMetadata(context.Background(), connect.NewRequest(&stencilv1beta1.GetSchemaMetadataRequest{
		NamespaceId: namespaceID,
		SchemaId:    schemaID,
	}))
	if err != nil {
		return nil, nil, err
	}

	return data, metaRes.Msg, nil
}
