package cmd

import (
	"context"
	"fmt"
	"os"

	"connectrpc.com/connect"
	"github.com/MakeNowJust/heredoc"
	"github.com/raystack/salt/cli/printer"
	stencilv1beta1 "github.com/raystack/stencil/gen/raystack/stencil/v1beta1"
	"github.com/spf13/cobra"
)

func checkSchemaCmd(cdk *CDK) *cobra.Command {
	var comp, file, namespaceID string

	cmd := &cobra.Command{
		Use:   "check <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Check schema compatibility",
		Long: heredoc.Doc(`
			Check schema compatibility of a local schema
			against a remote schema(against) on stencil server.`),
		Example: heredoc.Doc(`
			$ stencil schema check <id> -n raystack -c COMPATIBILITY_BACKWARD -F ./booking.desc
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			fileData, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			client, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}

			schemaID := args[0]

			req := &stencilv1beta1.CheckCompatibilityRequest{
				Data:          fileData,
				NamespaceId:   namespaceID,
				SchemaId:      schemaID,
				Compatibility: stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp]),
			}

			_, err = client.CheckCompatibility(context.Background(), connect.NewRequest(req))
			if err != nil {
				return err
			}

			spinner.Stop()
			fmt.Printf("\n%s Schema is compatible.\n", printer.Green(printer.Icon("success")))
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "Parent namespace ID")
	_ = cmd.MarkFlagRequired("namespace")

	cmd.Flags().StringVarP(&comp, "comp", "c", "", "Schema compatibility")
	_ = cmd.MarkFlagRequired("comp")

	cmd.Flags().StringVarP(&file, "file", "F", "", "Path to the schema file")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
