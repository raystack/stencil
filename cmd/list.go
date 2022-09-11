package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	"github.com/odpf/salt/term"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
)

func listSchemaCmd() *cobra.Command {
	var host, namespace string
	var req stencilv1beta1.ListSchemasRequest

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all schemas",
		Args:  cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			$ stencil schema list -n odpf
	    	`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			client, cancel, err := createClient(cmd)
			if err != nil {
				return err
			}
			defer cancel()

			req.Id = namespace
			res, err := client.ListSchemas(context.Background(), &req)
			if err != nil {
				return err
			}

			schemas := res.GetSchemas()

			// TODO(Ravi): List schemas should also handle namespace not found
			if len(schemas) == 0 {
				fmt.Printf("No schema found in namespace %s.\n", term.Blue(namespace))
				return nil
			}

			index := 1
			report := [][]string{}
			for _, s := range schemas {
				meta, _ := fetchMeta(client, namespace, s)
				c := meta.GetCompatibility().String()
				f := meta.GetFormat().String()
				a := meta.GetAuthority()

				if a == "" {
					a = "-"
				}

				report = append(report, []string{term.Greenf("#%02d", index), s, term.Grey(dict[f]), term.Grey(dict[c]), a})
				index++
			}

			spinner.Stop()
			fmt.Printf("\nShowing %d of %d schemas \n\n", len(schemas), len(schemas))
			printer.Table(os.Stdout, report)
			return nil
		},
	}

	// TODO(Ravi): Namespace should be optional.
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().StringVar(&host, "host", "", "Stencil host address eg: localhost:8000")

	return cmd
}
