package cmd

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	"github.com/odpf/salt/term"
	"github.com/spf13/cobra"
)

func downloadSchemaCmd() *cobra.Command {
	var host, output, namespaceID string
	var version int32
	var data []byte

	cmd := &cobra.Command{
		Use:   "download <id>",
		Short: "Download a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema download customer -n=odpf --version 1
	    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()
			client, cancel, err := createClient(cmd)
			if err != nil {
				return err
			}
			defer cancel()

			data, _, err = fetchSchemaAndMeta(client, version, namespaceID, args[0])
			if err != nil {
				return err
			}
			spinner.Stop()

			err = os.WriteFile(output, data, 0666)
			if err != nil {
				return err
			}

			fmt.Printf("%s Schema successfully written to %s\n", term.Green(term.SuccessIcon()), output)
			return nil
		},
	}
	cmd.Flags().StringVar(&host, "host", "", "Stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "Parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().Int32VarP(&version, "version", "v", 0, "Version of the schema")

	cmd.Flags().StringVarP(&output, "output", "o", "", "Path to the output file")
	cmd.MarkFlagRequired("output")

	return cmd
}
