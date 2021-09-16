package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/emicklei/dot"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func forEachMessage(msgs protoreflect.MessageDescriptors, f func(protoreflect.MessageDescriptor)) {
	for i := 0; i < msgs.Len(); i++ {
		msg := msgs.Get(i)
		f(msg)
		forEachMessage(msg.Messages(), f)
	}
}
func forEachField(fields protoreflect.FieldDescriptors, f func(protoreflect.FieldDescriptor)) {
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		f(field)
	}
}

// Visualize creates a new cobra command for visualize descriptor
func Visualize() *cobra.Command {

	var host, filePath string
	var req stencilv1.DownloadDescriptorRequest

	cmd := &cobra.Command{
		Use:   "visualize",
		Short: "Visualize filedescriptorset file",
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
			res, err := client.DownloadDescriptor(context.Background(), &req)
			if err != nil {
				return err
			}

			msg := &descriptorpb.FileDescriptorSet{}
			err = proto.Unmarshal(res.Data, msg)
			if err != nil {
				return fmt.Errorf("invalid file descriptorset file. %w", err)
			}

			files, err := protodesc.NewFiles(msg)
			if err != nil {
				return fmt.Errorf("file is not fully contained descriptor file.%w", err)
			}

			di := dot.NewGraph(dot.Directed)
			files.RangeFiles(func(file protoreflect.FileDescriptor) bool {

				subgraph := di.Subgraph(fmt.Sprintf("%s/%s", file.Package(), file.Path()), dot.ClusterOption{})
				forEachMessage(file.Messages(), func(msg protoreflect.MessageDescriptor) {
					group := subgraph.Subgraph(string(msg.FullName()), dot.ClusterOption{})
					forEachField(msg.Fields(), func(field protoreflect.FieldDescriptor) {
						node := group.Node(string(field.FullName()))
						node.Attr("label", string(field.FullName()))
						node.Attr("shape", "record")
						node.Attr("color", "0.650 0.700 0.700")
						node.Attr("ratio", "compress")
					})

				})
				return true
			})
			err = os.WriteFile(filePath, []byte(di.String()), 0666)
			return err
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
	cmd.Flags().StringVar(&filePath, "output", "", "write to file")
	cmd.Flags().StringSliceVar(&req.Fullnames, "fullnames", []string{}, "provide fully qualified proto full names. You can provide multiple names separated by \",\" Eg: google.protobuf.FileDescriptorProto,google.protobuf.FileDescriptorSet")
	return cmd
}
