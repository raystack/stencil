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
				childNode := di.Node(fmt.Sprintf("%s\n%s", string(file.Package()), file.Path()))
				childNode.Attr("shape", "note")
				childNode.Attr("style", "filled")
				childNode.Attr("fillcolor", "cornsilk")

				for i := 0; i < file.Imports().Len(); i++ {
					imp := file.Imports().Get(i)
					parentNode := di.Node(fmt.Sprintf("%s\n%s", string(imp.Package()), imp.Path()))
					parentNode.Attr("shape", "note")
					parentNode.Attr("style", "filled")
					parentNode.Attr("fillcolor", "cornsilk")
					di.Edge(childNode, parentNode, "depends on")
				}
				return true
			})
			err = os.WriteFile(filePath, []byte(di.String()), 0666)
			fmt.Println(".dot file has been created in", filePath)
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
	cmd.Flags().StringVar(&filePath, "output", "./proto_vis.dot", "write to file")
	cmd.Flags().StringSliceVar(&req.Fullnames, "fullnames", []string{}, "provide fully qualified proto full names. You can provide multiple names separated by \",\" Eg: google.protobuf.FileDescriptorProto,google.protobuf.FileDescriptorSet")
	return cmd
}
