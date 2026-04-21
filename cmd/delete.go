package cmd

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/MakeNowJust/heredoc"
	"github.com/raystack/salt/cli/printer"
	stencilv1beta1 "github.com/raystack/stencil/gen/raystack/stencil/v1beta1"
	"github.com/spf13/cobra"
)

func deleteSchemaCmd(cdk *CDK) *cobra.Command {
	var namespaceID string
	var version int32

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema delete booking -n raystack
	    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			client, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}

			schemaID := args[0]

			if version == 0 {
				_, err = client.DeleteSchema(context.Background(), connect.NewRequest(&stencilv1beta1.DeleteSchemaRequest{
					NamespaceId: namespaceID,
					SchemaId:    schemaID,
				}))
				if err != nil {
					return err
				}
			} else {
				_, err = client.DeleteVersion(context.Background(), connect.NewRequest(&stencilv1beta1.DeleteVersionRequest{
					NamespaceId: namespaceID,
					SchemaId:    schemaID,
					VersionId:   version,
				}))
				if err != nil {
					return err
				}
			}

			spinner.Stop()
			fmt.Printf("Schema successfully deleted")
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "Parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().Int32VarP(&version, "version", "v", 0, "Particular version to be deleted")

	return cmd
}
