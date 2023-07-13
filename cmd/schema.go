package cmd

import (
	"context"

	stencilv1beta1 "github.com/raystack/stencil/proto/raystack/stencil/v1beta1"
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
