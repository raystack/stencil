package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/raystack/salt/cli/printer"
	stencilv1beta1 "github.com/raystack/stencil/proto/raystack/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func infoSchemaCmd(cdk *CDK) *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "info <id>",
		Short: "View schema information",
		Long:  "Display the information about a schema.",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema info events -n raystack
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()
			client, cancel, err := createClient(cmd, cdk)
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
					fmt.Printf("%s Schema with id '%s' not found.\n", printer.Red(printer.Icon("failure")), args[0])
					return nil
				}
				return err
			}

			fmt.Printf("\n%s\n", printer.Blue(args[0]))
			fmt.Printf("\n%s\n\n", printer.Grey("No description provided"))
			fmt.Printf("%s \t %s \n", printer.Grey("Namespace:"), namespace)
			fmt.Printf("%s \t %s \n", printer.Grey("Format:"), dict[info.GetFormat().String()])
			fmt.Printf("%s \t %s \n", printer.Grey("Compatibility:"), dict[info.GetCompatibility().String()])
			fmt.Printf("%s \t %s \n\n", printer.Grey("Authority:"), dict[info.GetAuthority()])
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Provide schema namespace")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func versionSchemaCmd(cdk *CDK) *cobra.Command {
	var namespaceID string
	var req stencilv1beta1.ListVersionsRequest

	cmd := &cobra.Command{
		Use:   "version",
		Short: "View versions of a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema version booking -n raystack
	    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			client, cancel, err := createClient(cmd, cdk)
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
					printer.Greenf("#%v", strconv.FormatInt(int64(v), 10)),
					"-",
					"-",
				})
			}
			fmt.Printf("\nShowing %[1]d of %[1]d versions for %s\n \n", len(versions), schemaID)
			printer.Table(os.Stdout, report)
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	return cmd
}
