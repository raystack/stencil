package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/odpf/stencil/graph"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Graph creates a new cobra command for descriptor set dependencies graph
func GraphCmd() *cobra.Command {

	var host, filePath string
	var req stencilv1.DownloadDescriptorRequest

	cmd := &cobra.Command{
		Use:     "graph",
		Aliases: []string{"g"},
		Short:   "Generate file descriptorset dependencies graph",
		Args:    cobra.NoArgs,
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
			res, err := client.DownloadDescriptor(context.Background(), &req)
			if err != nil {
				return err
			}

			msg := &descriptorpb.FileDescriptorSet{}
			err = proto.Unmarshal(res.Data, msg)
			if err != nil {
				return fmt.Errorf("invalid file descriptorset file. %w", err)
			}

			graph, err := graph.GetProtoFileDependencyGraph(msg)
			if err != nil {
				return err
			}
			if err = os.WriteFile(filePath, []byte(graph.String()), 0666); err != nil {
				return err
			}

			fmt.Println(".dot file has been created in", filePath)
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")
	cmd.Flags().StringVar(&req.Namespace, "namespace", "", "provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")
	cmd.Flags().StringVar(&req.Name, "name", "", "provide proto repo name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&req.Version, "version", "", "provide semantic version compatible value")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringVar(&filePath, "output", "./proto_vis.dot", "write to file")
	cmd.Flags().StringSliceVar(&req.Fullnames, "fullnames", []string{}, "provide fully qualified proto full names. You can provide multiple names separated by \",\" Eg: google.protobuf.FileDescriptorProto,google.protobuf.FileDescriptorSet")
	return cmd
}
