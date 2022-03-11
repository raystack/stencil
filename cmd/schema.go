package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/odpf/salt/printer"
	"github.com/odpf/stencil/graph"
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
			$ stencil schema view
			$ stencil schema edit
			$ stencil schema delete
			$ stencil schema version
			$ stencil schema graph
			$ stencil schema print
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
	cmd.AddCommand(printCmd())
	cmd.AddCommand(graphCmd())

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
		Use:   "edit",
		Short: "Edit a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema edit <schema-id> --namespace=<namespace-id> --comp=<schema-compatibility>
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
	var resMetadata *stencilv1beta1.GetSchemaMetadataResponse

	cmd := &cobra.Command{
		Use:   "view",
		Short: "View a schema",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema view <schema-id> --namespace=<namespace-id> --version <version> --metadata <metadata>
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

			data, resMetadata, err = fetchSchemaAndMetadata(client, version, namespaceID, schemaID)
			if err != nil {
				return err
			}
			spinner.Stop()

			err = os.WriteFile(output, data, 0666)
			if err != nil {
				return err
			}

			fmt.Printf("Schema successfully written to %s\n", output)

			if resMetadata == nil || !metadata {
				return nil
			}

			report := [][]string{}

			fmt.Printf("\nMETADATA\n")
			report = append(report, []string{"FORMAT", "COMPATIBILITY", "AUTHORITY"})

			report = append(report, []string{
				stencilv1beta1.Schema_Format_name[int32(resMetadata.GetFormat())],
				stencilv1beta1.Schema_Compatibility_name[int32(resMetadata.GetCompatibility())],
				resMetadata.GetAuthority(),
			})

			printer.Table(os.Stdout, report)

			return nil
		},
	}
	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "parent namespace ID")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().Int32VarP(&version, "version", "v", 0, "version of the schema")

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
	var fullname string
	var host string
	var namespace string
	var earlierVersion int32
	var laterVersion int32

	var schemaFetcher = func(req *stencilv1beta1.GetSchemaRequest, client stencilv1beta1.StencilServiceClient) ([]byte, error) {
		res, err := client.GetSchema(context.Background(), req)
		if err != nil {
			return nil, err
		}
		return res.Data, nil
	}
	var protoSchemaFetcher = func(req *stencilv1beta1.GetSchemaRequest, client stencilv1beta1.StencilServiceClient) ([]byte, error) {
		if fullname == "" {
			return nil, fmt.Errorf("fullname flag is mandator for FORMAT_PROTO")
		}
		res, err := client.GetSchema(context.Background(), req)
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
			return nil, fmt.Errorf("fail to convert json. %w", err)
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
			defer spinner.Stop()

			schemaID := args[0]

			metaReq := stencilv1beta1.GetSchemaMetadataRequest{
				NamespaceId: namespace,
				SchemaId:    schemaID,
			}
			eReq := &stencilv1beta1.GetSchemaRequest{
				NamespaceId: namespace,
				SchemaId:    schemaID,
				VersionId:   earlierVersion,
			}
			lReq := &stencilv1beta1.GetSchemaRequest{
				NamespaceId: namespace,
				SchemaId:    schemaID,
				VersionId:   laterVersion,
			}

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			client := stencilv1beta1.NewStencilServiceClient(conn)

			meta, err := client.GetSchemaMetadata(context.Background(), &metaReq)
			if err != nil {
				return err
			}

			var getSchema = schemaFetcher
			if meta.Format == *stencilv1beta1.Schema_FORMAT_PROTOBUF.Enum() {
				getSchema = protoSchemaFetcher
			}

			eJson, err := getSchema(eReq, client)
			if err != nil {
				return err
			}

			lJson, err := getSchema(lReq, client)
			if err != nil {
				return err
			}

			d, err := gojsondiff.New().Compare(eJson, lJson)
			if err != nil {
				return err
			}

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
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "parent namespace ID")
	cmd.MarkFlagRequired("namespace")
	cmd.Flags().Int32Var(&earlierVersion, "earlier-version", 0, "earlier version of the schema")
	cmd.MarkFlagRequired("earlier-version")
	cmd.Flags().Int32Var(&laterVersion, "later-version", 0, "later version of the schema")
	cmd.MarkFlagRequired("later-version")
	cmd.Flags().StringVar(&fullname, "fullname", "", "only required for FORMAT_PROTO. fullname of proto schema eg: odpf.common.v1.Version")
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

func printCmd() *cobra.Command {
	var output, filterPathPrefix, host, namespaceID, schemaID string
	var version int32

	cmd := &cobra.Command{
		Use:   "print",
		Short: "Prints snapshot details into .proto files",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema print <schema-id> --namespace=<namespace-id> --version <version> --output=<output-path> --filter-path=<path-prefix>
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			client := stencilv1beta1.NewStencilServiceClient(conn)

			schemaID := args[0]

			data, resMetadata, err := fetchSchemaAndMetadata(client, version, namespaceID, schemaID)
			if err != nil {
				return err
			}

			format := stencilv1beta1.Schema_Format_name[int32(resMetadata.GetFormat())]

			if format == "FORMAT_AVRO" || format == "FORMAT_JSON" {
				if output == "" {
					fmt.Printf("\n// ----\n// SCHEMA\n// ----\n\n")
					_, err := os.Stdout.Write(data)
					if err != nil {
						return fmt.Errorf("schema is not valid. %w", err)
					}
				} else {
					err = os.WriteFile(output, data, 0666)
					if err != nil {
						return err
					}

					fmt.Printf("Schema successfully written to %s\n", output)
				}
			} else {
				fds := &descriptorpb.FileDescriptorSet{}
				if err := proto.Unmarshal(data, fds); err != nil {
					return fmt.Errorf("descriptor set file is not valid. %w", err)
				}
				fdsMap, err := desc.CreateFileDescriptorsFromSet(fds)
				if err != nil {
					return err
				}

				var filteredFds []*desc.FileDescriptor
				for fdName, fd := range fdsMap {
					if filterPathPrefix != "" && !strings.HasPrefix(fdName, filterPathPrefix) {
						continue
					}
					filteredFds = append(filteredFds, fd)
				}

				protoPrinter := &protoprint.Printer{}

				if output == "" {
					for _, fd := range filteredFds {
						protoAsString, err := protoPrinter.PrintProtoToString(fd)
						if err != nil {
							return err
						}
						fmt.Printf("\n// ----\n// %s\n// ----\n%s", fd.GetName(), protoAsString)
					}
				} else {
					if err := protoPrinter.PrintProtosToFileSystem(filteredFds, output); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().StringVarP(&schemaID, "schema", "s", "", "provide proto repo name")
	cmd.MarkFlagRequired("name")

	cmd.Flags().Int32VarP(&version, "version", "v", 0, "provide version number")

	cmd.Flags().StringVarP(&output, "output", "o", "", "the directory path to write the descriptor files, default is to print on stdout")

	cmd.Flags().StringVar(&filterPathPrefix, "filter-path", "", "filter protocol buffer files by path prefix, e.g., --filter-path=google/protobuf")

	return cmd
}

func graphCmd() *cobra.Command {
	var host, output, namespaceID string
	var version int32

	cmd := &cobra.Command{
		Use:     "graph",
		Aliases: []string{"g"},
		Short:   "Generate file descriptorset dependencies graph",
		Args:    cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema graph <schema-id> --namespace=<namespace-id> --version=<version> --output=<output-path>
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			client := stencilv1beta1.NewStencilServiceClient(conn)

			schemaID := args[0]

			data, resMetadata, err := fetchSchemaAndMetadata(client, version, namespaceID, schemaID)
			if err != nil {
				return err
			}

			format := stencilv1beta1.Schema_Format_name[int32(resMetadata.GetFormat())]
			if format != "FORMAT_PROTOBUF" {
				fmt.Printf("cannot create graph for %s", format)
				return nil
			}

			msg := &descriptorpb.FileDescriptorSet{}
			err = proto.Unmarshal(data, msg)
			if err != nil {
				return fmt.Errorf("invalid file descriptorset file. %w", err)
			}

			graph, err := graph.GetProtoFileDependencyGraph(msg)
			if err != nil {
				return err
			}
			if err = os.WriteFile(output, []byte(graph.String()), 0666); err != nil {
				return err
			}

			fmt.Println(".dot file has been created in", output)
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().Int32VarP(&version, "version", "v", 0, "provide version number")

	cmd.Flags().StringVarP(&output, "output", "o", "./proto_vis.dot", "write to .dot file")

	return cmd
}

func fetchSchemaAndMetadata(client stencilv1beta1.StencilServiceClient, version int32, namespaceID, schemaID string) ([]byte, *stencilv1beta1.GetSchemaMetadataResponse, error) {
	var req stencilv1beta1.GetSchemaRequest
	var reqLatest stencilv1beta1.GetLatestSchemaRequest
	var reqMetadata stencilv1beta1.GetSchemaMetadataRequest
	var data []byte

	if version != 0 {
		req.NamespaceId = namespaceID
		req.SchemaId = schemaID
		req.VersionId = version
		res, err := client.GetSchema(context.Background(), &req)
		if err != nil {
			return nil, nil, err
		}
		data = res.GetData()
	} else {
		reqLatest.NamespaceId = namespaceID
		reqLatest.SchemaId = schemaID
		res, err := client.GetLatestSchema(context.Background(), &reqLatest)
		if err != nil {
			return nil, nil, err
		}
		data = res.GetData()
	}

	reqMetadata.NamespaceId = namespaceID
	reqMetadata.SchemaId = schemaID
	resMetadata, err := client.GetSchemaMetadata(context.Background(), &reqMetadata)
	if err != nil {
		return data, nil, err
	}

	return data, resMetadata, nil
}
