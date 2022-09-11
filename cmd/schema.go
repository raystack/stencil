package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	"github.com/odpf/salt/term"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

func SchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "schema",
		Aliases: []string{"schemas"},
		Short:   "Manage schemas",
		Long: heredoc.Doc(`
			Work with schemas.
		`),
		Annotations: map[string]string{
			"group": "core",
		},
	}

	cmd.AddCommand(listSchemaCmd())
	cmd.AddCommand(infoSchemaCmd())
	cmd.AddCommand(printSchemaCmd())
	cmd.AddCommand(downloadSchemaCmd())

	cmd.AddCommand(createSchemaCmd())
	cmd.AddCommand(checkSchemaCmd())
	cmd.AddCommand(updateSchemaCmd())
	cmd.AddCommand(deleteSchemaCmd())
	cmd.AddCommand(diffSchemaCmd())
	cmd.AddCommand(versionSchemaCmd())

	cmd.AddCommand(graphCmd())

	return cmd
}

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
			spinner.Stop()

			// TODO(Ravi): List schemas should also handle namespace not found
			if len(schemas) == 0 {
				fmt.Printf("No schema found in namespace %s.\n", term.Blue(namespace))
				return nil
			}

			fmt.Printf("\nShowing %d of %d schemas \n\n", len(schemas), len(schemas))
			index := 1

			report := [][]string{}
			for _, s := range schemas {
				report = append(report, []string{term.Greenf("#%02d", index), s})
				index++
			}
			printer.Table(os.Stdout, report)
			return nil
		},
	}

	// TODO(Ravi): Namespace should be optional.
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")

	return cmd
}

func createSchemaCmd() *cobra.Command {
	var host, format, comp, filePath, namespaceID string
	var req stencilv1beta1.CreateSchemaRequest

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema create <schema-id> --namespace=<namespace-id> --format=<schema-format> –-comp=<schema-compatibility> –-filePath=<schema-filePath> 
	    	`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			fileData, err := ioutil.ReadFile(filePath)
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
	cmd.MarkFlagRequired("format")

	cmd.Flags().StringVarP(&comp, "comp", "c", "", "schema compatibility")
	cmd.MarkFlagRequired("comp")

	cmd.Flags().StringVarP(&filePath, "filePath", "F", "", "path to the schema file")
	cmd.MarkFlagRequired("filePath")

	return cmd
}

func updateSchemaCmd() *cobra.Command {
	var host, comp, namespaceID string
	var req stencilv1beta1.UpdateSchemaMetadataRequest

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema edit <schema-id> --namespace=<namespace-id> --comp=<schema-compatibility>
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
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp])

			_, err = client.UpdateSchemaMetadata(context.Background(), &req)
			if err != nil {
				return err
			}

			spinner.Stop()

			fmt.Printf("Schema successfully updated")
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().StringVarP(&comp, "comp", "c", "", "schema compatibility")
	cmd.MarkFlagRequired("comp")

	return cmd
}

func deleteSchemaCmd() *cobra.Command {
	var host, namespaceID string
	var req stencilv1beta1.DeleteSchemaRequest
	var reqVer stencilv1beta1.DeleteVersionRequest
	var version int32

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema delete <schema-id> --namespace=<namespace-id>
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

			if version == 0 {
				req.NamespaceId = namespaceID
				req.SchemaId = schemaID

				_, err = client.DeleteSchema(context.Background(), &req)
				if err != nil {
					return err
				}
			} else {
				reqVer.NamespaceId = namespaceID
				reqVer.SchemaId = schemaID
				reqVer.VersionId = version

				_, err = client.DeleteVersion(context.Background(), &reqVer)
				if err != nil {
					return err
				}
			}

			spinner.Stop()

			fmt.Printf("schema successfully deleted")

			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().Int32VarP(&version, "version", "v", 0, "particular version to be deleted")

	return cmd
}

func fetchSchemaAndMeta(client stencilv1beta1.StencilServiceClient, version int32, namespaceID, schemaID string) ([]byte, *stencilv1beta1.GetSchemaMetadataResponse, error) {
	var req stencilv1beta1.GetSchemaRequest
	var reqLatest stencilv1beta1.GetLatestSchemaRequest
	var data []byte

	ctx := context.Background()

	if version != 0 {
		req.NamespaceId = namespaceID
		req.SchemaId = schemaID
		req.VersionId = version
		res, err := client.GetSchema(ctx, &req)
		if err != nil {
			return nil, nil, err
		}
		data = res.GetData()
	} else {
		reqLatest.NamespaceId = namespaceID
		reqLatest.SchemaId = schemaID
		res, err := client.GetLatestSchema(ctx, &reqLatest)
		if err != nil {
			return nil, nil, err
		}
		data = res.GetData()
	}

	reqMeta := stencilv1beta1.GetSchemaMetadataRequest{
		NamespaceId: namespaceID,
		SchemaId:    schemaID,
	}
	meta, err := client.GetSchemaMetadata(context.Background(), &reqMeta)

	if err != nil {
		return nil, nil, err
	}

	return data, meta, nil
}

// func fetchMeta(client stencilv1beta1.StencilServiceClient, namespace string, schema string) (*stencilv1beta1.GetSchemaMetadataResponse, error) {
// 	req := stencilv1beta1.GetSchemaMetadataRequest{
// 		NamespaceId: namespace,
// 		SchemaId:    schema,
// 	}
// 	return client.GetSchemaMetadata(context.Background(), &req)
// }
