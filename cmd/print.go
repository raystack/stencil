package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// PrintCmd creates a new cobra command for getting .proto files
func PrintCmd() *cobra.Command {
	var (
		req              stencilv1.GetSchemaRequest
		pathDir          string
		filterPathPrefix string
		host             string
	)
	cmd := &cobra.Command{
		Use:   "print",
		Short: "prints snapshot details into .proto files",
		Args:  cobra.NoArgs,
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			client := stencilv1.NewStencilServiceClient(conn)
			res, err := client.GetSchema(context.Background(), &req)
			if err != nil {
				return err
			}

			fds := &descriptorpb.FileDescriptorSet{}
			if err := proto.Unmarshal(res.GetData(), fds); err != nil {
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
			if pathDir == "" {
				for _, fd := range filteredFds {
					protoAsString, err := protoPrinter.PrintProtoToString(fd)
					if err != nil {
						return err
					}
					fmt.Printf("\n// ----\n// %s\n// ----\n%s", fd.GetName(), protoAsString)
				}
			} else {
				if err := protoPrinter.PrintProtosToFileSystem(filteredFds, pathDir); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")
	cmd.Flags().StringVar(&req.NamespaceId, "namespace", "", "provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")
	cmd.Flags().StringVar(&req.SchemaId, "name", "", "provide proto repo name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().Int32Var(&req.VersionId, "version", 0, "provide version number")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringVar(&pathDir, "output", "", "the directory path to write the descriptor files, default is to print on stdout")
	cmd.Flags().StringVar(&filterPathPrefix, "filter-path", "", "filter protocol buffer files by path prefix, e.g., --filter-path=google/protobuf")
	return cmd
}
