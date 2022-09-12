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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func infoSchemaCmd() *cobra.Command {
	var host, namespace string

	cmd := &cobra.Command{
		Use:   "info <id>",
		Short: "View schema information",
		Long:  "Display the information about a schema.",
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
			spinner.Stop()
			if err != nil {
				errStatus, _ := status.FromError(err)
				if codes.NotFound == errStatus.Code() {
					fmt.Printf("%s Schema with id '%s' not found.\n", term.Red(term.FailureIcon()), args[0])
					return nil
				}
				return err
			}

			fmt.Printf("\n%s\n", term.Blue(args[0]))
			fmt.Printf("\n%s\n\n", term.Grey("No description provided"))
			fmt.Printf("%s \t %s \n", term.Grey("Namespace:"), namespace)
			fmt.Printf("%s \t %s \n", term.Grey("Format:"), dict[info.GetFormat().String()])
			fmt.Printf("%s \t %s \n", term.Grey("Compatibility:"), dict[info.GetCompatibility().String()])
			fmt.Printf("%s \t %s \n\n", term.Grey("Authority:"), dict[info.GetAuthority()])
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
		Short: "View versions of a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema version booking -n odpf
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

			versions := res.GetVersions()
			spinner.Stop()

			if len(versions) == 0 {
				fmt.Printf("No version found for %s in %s", schemaID, namespaceID)
				return nil
			}

			report := [][]string{}
			report = append(report, []string{"VERSION", "CREATED", "MESSAGE"})

			for _, v := range versions {
				report = append(report, []string{
					term.Greenf("#%v", strconv.FormatInt(int64(v), 10)),
					"-",
					"-",
				})
			}
			fmt.Printf("\nShowing %[1]d of %[1]d versions for %s\n \n", len(versions), schemaID)
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
