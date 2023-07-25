package cmd

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/printer"
	stencilv1beta1 "github.com/goto/stencil/proto/v1beta1"
	"github.com/spf13/cobra"
)

func editSchemaCmd(cdk *CDK) *cobra.Command {
	var comp, namespaceID string
	var req stencilv1beta1.UpdateSchemaMetadataRequest

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema edit booking -n goto -c COMPATIBILITY_BACKWARD
	    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			client, cancel, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}
			defer cancel()

			schemaID := args[0]

			req.NamespaceId = namespaceID
			req.SchemaId = schemaID
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp])

			_, err = client.UpdateSchemaMetadata(context.Background(), &req)
			if err != nil {
				return err
			}

			spinner.Stop()
			fmt.Printf("Schema successfully updated")
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "Parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().StringVarP(&comp, "comp", "c", "", "Schema compatibility")
	cmd.MarkFlagRequired("comp")

	return cmd
}
