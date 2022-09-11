package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	"github.com/odpf/salt/prompt"
	"github.com/odpf/salt/term"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

func createSchemaCmd() *cobra.Command {
	var host, format, comp, filePath, namespaceID string
	var req stencilv1beta1.CreateSchemaRequest

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema create booking -n odpf –F booking.json
			$ stencil schema create booking -n odpf -f FORMAT_JSON –c COMPATABILITY_BACKWARD –F ./booking.json 
	    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompter := prompt.New()
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			req.Data = fileData

			if format == "" {
				formatAnswer, _ := prompter.Select("Select schema format:", formats[0], formats)
				format = formats[formatAnswer]
			}

			if comp == "" {
				formatAnswer, _ := prompter.Select("Select schema compatibility:", comps[0], comps)
				fmt.Println()
				comp = comps[formatAnswer]
			}

			spinner := printer.Spin("")
			defer spinner.Stop()
			client, cancel, err := createClient(cmd)
			if err != nil {
				return err
			}
			defer cancel()

			schemaID := args[0]

			req.NamespaceId = namespaceID
			req.SchemaId = schemaID
			req.Format = stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[format])
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp])

			res, err := client.CreateSchema(context.Background(), &req)
			if err != nil {
				errStatus := status.Convert(err)
				return errors.New(errStatus.Message())
			}

			id := res.GetId()

			spinner.Stop()
			fmt.Printf("\n%s Created schema with id %s.\n", term.Green(term.SuccessIcon()), term.Cyan(id))
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "Namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().StringVarP(&format, "format", "f", "", "schema format")

	cmd.Flags().StringVarP(&comp, "comp", "c", "", "schema compatibility")

	cmd.Flags().StringVarP(&filePath, "filePath", "F", "", "path to the schema file")
	cmd.MarkFlagRequired("filePath")

	return cmd
}
