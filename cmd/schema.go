package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	stencilv1beta1 "github.com/odpf/stencil/server/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func SchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "schema",
		Aliases: []string{"schema"},
		Short:   "Manage schema",
		Long: heredoc.Doc(`
			Work with schemas.
		`),
		Example: heredoc.Doc(`
			$ stencil schema list
			$ stencil schema create
			$ stencil schema get
			$ stencil schema update
			$ stencil schema delete
			$ stencil schema version
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
	}

	cmd.AddCommand(createSchemaCmd())
	cmd.AddCommand(listSchemaCmd())
	cmd.AddCommand(getSchemaCmd())
	cmd.AddCommand(updateSchemaCmd())
	cmd.AddCommand(deleteSchemaCmd())
	cmd.AddCommand(diffSchemaCmd())
	cmd.AddCommand(versionSchemaCmd())

	return cmd
}

func listSchemaCmd() *cobra.Command {
	var host string
	var req stencilv1beta1.ListSchemasRequest

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all schemas",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema list <namespace-id>
	    	`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			namespaceID := args[0]
			req.Id = namespaceID

			client := stencilv1beta1.NewStencilServiceClient(conn)
			res, err := client.ListSchemas(context.Background(), &req)
			if err != nil {
				return err
			}

			report := [][]string{}

			schemas := res.GetSchemas()

			spinner.Stop()

			if len(schemas) == 0 {
				fmt.Printf("%s has no schemas", namespaceID)
				return nil
			}

			fmt.Printf(" \nShowing %d schemas \n", len(schemas))

			report = append(report, []string{"SCHEMA"})

			for _, s := range schemas {
				report = append(report, []string{
					s,
				})
			}
			printer.Table(os.Stdout, report)
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

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
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			fileData, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}
			req.Data = fileData

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			schemaID := args[0]

			req.NamespaceId = namespaceID
			req.SchemaId = schemaID
			req.Format = stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[format])
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp])

			client := stencilv1beta1.NewStencilServiceClient(conn)
			res, err := client.CreateSchema(context.Background(), &req)
			if err != nil {
				errStatus := status.Convert(err)
				return errors.New(errStatus.Message())
			}

			id := res.GetId()

			spinner.Stop()
			fmt.Printf("schema successfully created with id: %s", id)
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "parent namespace ID")
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
		Use:   "update",
		Short: "Edit a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema update <schema-id> --namespace=<namespace-id> –-comp=<schema-compatibility>
	    	`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			schemaID := args[0]

			req.NamespaceId = namespaceID
			req.SchemaId = schemaID
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp])

			client := stencilv1beta1.NewStencilServiceClient(conn)
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

func getSchemaCmd() *cobra.Command {
	var host, output, namespaceID string
	var version int32
	var metadata bool
	var data []byte
	var req stencilv1beta1.GetSchemaRequest
	var reqLatest stencilv1beta1.GetLatestSchemaRequest
	var reqMetadata stencilv1beta1.GetSchemaMetadataRequest
	var resMetadata *stencilv1beta1.GetSchemaMetadataResponse

	cmd := &cobra.Command{
		Use:   "get",
		Short: "View a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema get <schema-id> --namespace=<namespace-id> --version <version> --metadata <metadata>
	    	`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			schemaID := args[0]

			client := stencilv1beta1.NewStencilServiceClient(conn)

			if version != 0 {
				req.NamespaceId = namespaceID
				req.SchemaId = schemaID
				req.VersionId = version
				res, err := client.GetSchema(context.Background(), &req)
				if err != nil {
					return err
				}
				data = res.GetData()
			} else {
				reqLatest.NamespaceId = namespaceID
				reqLatest.SchemaId = schemaID
				res, err := client.GetLatestSchema(context.Background(), &reqLatest)
				if err != nil {
					return err
				}
				data = res.GetData()
			}

			if metadata {
				reqMetadata.NamespaceId = namespaceID
				reqMetadata.SchemaId = schemaID
				resMetadata, err = client.GetSchemaMetadata(context.Background(), &reqMetadata)
				if err != nil {
					return err
				}
			}

			spinner.Stop()

			err = os.WriteFile(output, data, 0666)
			if err != nil {
				return err
			}

			fmt.Printf("Schema successfully written to %s\n", output)

			if resMetadata != nil {
				report := [][]string{}

				fmt.Printf("\nMETADATA\n")
				report = append(report, []string{"FORMAT", "COMPATIBILITY", "AUTHORITY"})

				report = append(report, []string{
					stencilv1beta1.Schema_Format_name[int32(resMetadata.GetFormat())],
					stencilv1beta1.Schema_Compatibility_name[int32(resMetadata.GetCompatibility())],
					resMetadata.GetAuthority(),
				})

				printer.Table(os.Stdout, report)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().Int32VarP(&version, "version", "v", 0, "version of the schema")
	cmd.MarkFlagRequired("version")

	cmd.Flags().BoolVarP(&metadata, "metadata", "m", false, "set this flag to get metadata")
	cmd.MarkFlagRequired("metadata")

	cmd.Flags().StringVarP(&output, "output", "o", "", "path to the output file")
	cmd.MarkFlagRequired("output")

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
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			schemaID := args[0]

			client := stencilv1beta1.NewStencilServiceClient(conn)

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

func diffSchemaCmd() *cobra.Command {
	var namespace string
	var earlierVersion int32
	var laterVersion int32
	var host string
	var fullname string

	var getJsonMessage = func(req stencilv1beta1.GetSchemaRequest, conn *grpc.ClientConn) ([]byte, error) {
		client := stencilv1beta1.NewStencilServiceClient(conn)
		res, err := client.GetSchema(context.Background(), &req)
		if err != nil {
			return nil, err
		}
		fds := &descriptorpb.FileDescriptorSet{}
		if err := proto.Unmarshal(res.Data, fds); err != nil {
			return nil, fmt.Errorf("descriptor set file is not valid. %w", err)
		}
		files, err := protodesc.NewFiles(fds)
		if err != nil {
			return nil, fmt.Errorf("file is not fully contained descriptor file. hint: generate file descriptorset with --include_imports option. %w", err)
		}
		desc, err := files.FindDescriptorByName(protoreflect.FullName(fullname))
		if err != nil {
			return nil, fmt.Errorf("unable to find message. %w", err)
		}
		mDesc, ok := desc.(protoreflect.MessageDescriptor)
		if !ok {
			return nil, fmt.Errorf("not a message desc")
		}
		jsonByte, err := protojson.Marshal(protodesc.ToDescriptorProto(mDesc))
		if err != nil {
			return nil, fmt.Errorf("fail to convert o json. %w", err)
		}
		return jsonByte, nil
	}

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Diff(s) of two schema versions",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
		$ stencil schema diff <schema-id> --namespace=<namespace-id> --later-version=<later-version> --earlier-version=<earlier-version> --fullname=<fullname>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			schemaID := args[0]
			eReq := stencilv1beta1.GetSchemaRequest{
				NamespaceId: namespace,
				SchemaId:    schemaID,
				VersionId:   earlierVersion,
			}
			lReq := stencilv1beta1.GetSchemaRequest{
				NamespaceId: namespace,
				SchemaId:    schemaID,
				VersionId:   laterVersion,
			}

			host, _ := cmd.Flags().GetString("host")
			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			eJson, err := getJsonMessage(eReq, conn)
			if err != nil {
				return err
			}

			lJson, err := getJsonMessage(lReq, conn)
			if err != nil {
				return err
			}

			d, err := gojsondiff.New().Compare(eJson, lJson)

			var placeholder map[string]interface{}
			json.Unmarshal(eJson, &placeholder)
			config := formatter.AsciiFormatterConfig{
				ShowArrayIndex: true,
				Coloring:       true,
			}

			formatter := formatter.NewAsciiFormatter(placeholder, config)
			diffString, err := formatter.Format(d)
			if err != nil {
				return err
			}

			spinner.Stop()
			if !d.Modified() {
				fmt.Print("No diff!")
				return nil
			}
			fmt.Print(diffString)

			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")
	cmd.Flags().StringVar(&namespace, "namespace", "", "parent namespace ID")
	cmd.MarkFlagRequired("namespace")
	cmd.Flags().Int32Var(&earlierVersion, "earlier-version", 0, "earlier version of the schema")
	cmd.MarkFlagRequired("earlier-version")
	cmd.Flags().Int32Var(&laterVersion, "later-version", 0, "later version of the schema")
	cmd.MarkFlagRequired("later-version")
	cmd.Flags().StringVar(&fullname, "fullname", "", "fullname of proto schema eg: odpf.common.v1.Version")
	cmd.MarkFlagRequired("fullname")
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
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			schemaID := args[0]

			req.NamespaceId = namespaceID
			req.SchemaId = schemaID

			client := stencilv1beta1.NewStencilServiceClient(conn)
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
