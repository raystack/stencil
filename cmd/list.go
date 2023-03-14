package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	stencilv1beta1 "github.com/goto/stencil/proto/v1beta1"
	"github.com/odpf/salt/printer"
	"github.com/odpf/salt/term"
	"github.com/spf13/cobra"
)

func listSchemaCmd(cdk *CDK) *cobra.Command {
	var namespace string
	var req stencilv1beta1.ListSchemasRequest

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all schemas",
		Long: heredoc.Doc(`
			List schemas in a namespace.
		`),
		Args: cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			$ stencil schema list -n goto
	    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			client, cancel, err := createClient(cmd, cdk)
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
				spinner.Stop()
				fmt.Printf("No schema found in namespace %s\n", namespace)
				return nil
			}

			report := [][]string{}
			index := 1
			report = append(report, []string{
				term.Bold("INDEX"),
				term.Bold("NAME"),
				term.Bold("FORMAT"),
				term.Bold("COMPATIBILITY"),
				term.Bold("AUTHORITY"),
			})
			for _, s := range schemas {
				c := s.GetCompatibility().String()
				f := s.GetFormat().String()
				a := s.GetAuthority()

				if a == "" {
					a = "-"
				}
				report = append(report, []string{term.Greenf("#%d", index), s.GetName(), dict[f], dict[c], a})
				index++
			}

			spinner.Stop()
			fmt.Printf("\nShowing %d of %d schemas in %s\n\n", len(schemas), len(schemas), namespace)
			printer.Table(os.Stdout, report)
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace ID")
	cmd.MarkFlagRequired("namespace")

	return cmd
}
