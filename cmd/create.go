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

func createSchemaCmd(cdk *CDK) *cobra.Command {
	var format, comp, file, namespaceID string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema create booking -n raystack -F booking.json
			$ stencil schema create booking -n raystack -f FORMAT_JSON -c COMPATIBILITY_BACKWARD -F ./booking.json
	    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileData, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			spinner := printer.Spin("")
			defer spinner.Stop()
			client, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}

			schemaID := args[0]
			req := &stencilv1beta1.CreateSchemaRequest{
				NamespaceId:   namespaceID,
				SchemaId:      schemaID,
				Data:          fileData,
				Format:        stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[format]),
				Compatibility: stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp]),
			}

			res, err := client.CreateSchema(context.Background(), connect.NewRequest(req))
			if err != nil {
				connectErr, ok := err.(*connect.Error)
				if ok && connectErr.Code() == connect.CodeAlreadyExists {
					fmt.Printf("\n%s Schema with id '%s' already exist.\n", printer.Icon("failure"), args[0])
					return nil
				}
				return err
			}

			id := res.Msg.GetId()

			spinner.Stop()
			fmt.Printf("\n%s Created schema with id %s.\n", printer.Green(printer.Icon("success")), printer.Cyan(id))
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "Namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().StringVarP(&format, "format", "f", "", "Schema format")

	cmd.Flags().StringVarP(&comp, "comp", "c", "", "Schema compatibility")

	cmd.Flags().StringVarP(&file, "file", "F", "", "Path to the schema file")
	cmd.MarkFlagRequired("file")

	return cmd
}
