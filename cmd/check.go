package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	"github.com/odpf/salt/term"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

func checkSchemaCmd() *cobra.Command {
	var host, comp, filePath, namespaceID string
	var req stencilv1beta1.CheckCompatibilityRequest

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check schema compatibility",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema check <schema-id> --namespace=<namespace-id> comp=<schema-compatibility> filePath=<schema-filePath>
	    	`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			req.Data = fileData

			client, cancel, err := createClient(cmd)
			if err != nil {
				return err
			}
			defer cancel()

			schemaID := args[0]

			req.NamespaceId = namespaceID
			req.SchemaId = schemaID
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp])

			_, err = client.CheckCompatibility(context.Background(), &req)
			if err != nil {
				errStatus := status.Convert(err)
				return errors.New(errStatus.Message())
			}

			spinner.Stop()
			fmt.Println("schema is compatible")
			fmt.Printf("\n%s Schema is compatible.\n", term.Green(term.SuccessIcon()))
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().StringVarP(&comp, "comp", "c", "", "schema compatibility")
	cmd.MarkFlagRequired("comp")

	cmd.Flags().StringVarP(&filePath, "filePath", "F", "", "path to the schema file")
	cmd.MarkFlagRequired("filePath")

	return cmd
}
