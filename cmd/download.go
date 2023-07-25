package cmd

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/printer"
	"github.com/goto/salt/term"
	"github.com/spf13/cobra"
)

func downloadSchemaCmd(cdk *CDK) *cobra.Command {
	var output, namespaceID string
	var version int32
	var data []byte

	cmd := &cobra.Command{
		Use:   "download <id>",
		Short: "Download a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema download customer -n=goto --version 1
	    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()
			client, cancel, err := createClient(cmd, cdk)
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

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "Parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().Int32VarP(&version, "version", "v", 0, "Version of the schema")

	cmd.Flags().StringVarP(&output, "output", "o", "", "Path to the output file")
	cmd.MarkFlagRequired("output")

	return cmd
}
