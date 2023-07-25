package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/printer"
	"github.com/goto/salt/term"
	stencilv1beta1 "github.com/goto/stencil/proto/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

func checkSchemaCmd(cdk *CDK) *cobra.Command {
	var comp, file, namespaceID string
	var req stencilv1beta1.CheckCompatibilityRequest

	cmd := &cobra.Command{
		Use:   "check <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Check schema compatibility",
		Long: heredoc.Doc(`
			Check schema compatibility of a local schema
			against a remote schema(against) on stencil server.`),
		Example: heredoc.Doc(`
			$ stencil schema check <id> -n goto -c COMPATIBILITY_BACKWARD -F ./booking.desc
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			fileData, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			client, cancel, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}
			defer cancel()

			schemaID := args[0]

			req.Data = fileData
			req.NamespaceId = namespaceID
			req.SchemaId = schemaID
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp])

			_, err = client.CheckCompatibility(context.Background(), &req)
			if err != nil {
				errStatus := status.Convert(err)
				return errors.New(errStatus.Message())
			}

			spinner.Stop()
			fmt.Printf("\n%s Schema is compatible.\n", term.Green(term.SuccessIcon()))
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "Parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().StringVarP(&comp, "comp", "c", "", "Schema compatibility")
	cmd.MarkFlagRequired("comp")

	cmd.Flags().StringVarP(&file, "file", "F", "", "Path to the schema file")
	cmd.MarkFlagRequired("file")

	return cmd
}
