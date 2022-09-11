package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	"github.com/odpf/salt/term"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
)

func infoSchemaCmd() *cobra.Command {
	var host, namespace string

	cmd := &cobra.Command{
		Use:   "info <id>",
		Short: "Print a given schema snapshot",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema info events -n odpf
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			client, cancel, err := createClient(cmd)
			if err != nil {
				return err
			}
			defer cancel()

			req := stencilv1beta1.GetSchemaMetadataRequest{
				NamespaceId: namespace,
				SchemaId:    args[0],
			}

			info, err := client.GetSchemaMetadata(cmd.Context(), &req)
			if err != nil {
				return err
			}
			spinner.Stop()

			fmt.Printf("%s \t\t %s \n", term.Bold("Format:"), info.GetFormat())
			fmt.Printf("%s \t\t %s \n", term.Bold("Compatibility:"), info.GetCompatibility())
			fmt.Printf("%s \t\t %s \n", term.Bold("Authority:"), info.GetAuthority())
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "Stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Provide schema namespace")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func versionSchemaCmd() *cobra.Command {
	var host, namespaceID string
	var req stencilv1beta1.ListVersionsRequest

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Version(s) of a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema version <schema-id> --namespace=<namespace-id>
	    	`),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			res, err := client.ListVersions(context.Background(), &req)
			if err != nil {
				return err
			}

			report := [][]string{}
			versions := res.GetVersions()

			spinner.Stop()

			if len(versions) == 0 {
				fmt.Printf("%s has no versions in %s", schemaID, namespaceID)
				return nil
			}

			report = append(report, []string{"VERSIONS(s)"})

			for _, v := range versions {
				report = append(report, []string{
					strconv.FormatInt(int64(v), 10),
				})
			}
			printer.Table(os.Stdout, report)
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	return cmd
}
