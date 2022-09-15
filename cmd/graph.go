package cmd

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/stencil/pkg/graph"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func graphSchemaCmd(cdk *CDK) *cobra.Command {
	var output, namespaceID string
	var version int32

	cmd := &cobra.Command{
		Use:     "graph",
		Aliases: []string{"g"},
		Short:   "View schema dependencies graph",
		Args:    cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema graph booking -n odpf -v 1 -o ./vis.dot
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, cancel, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}
			defer cancel()

			schemaID := args[0]

			data, resMetadata, err := fetchSchemaAndMeta(client, version, namespaceID, schemaID)
			if err != nil {
				return err
			}

			format := stencilv1beta1.Schema_Format_name[int32(resMetadata.GetFormat())]
			if format != "FORMAT_PROTOBUF" {
				fmt.Printf("Graph is not supported for %s", format)
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

			fmt.Println("Created graph file at", output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().Int32VarP(&version, "version", "v", 0, "provide version number")

	cmd.Flags().StringVarP(&output, "output", "o", "./proto_vis.dot", "write to .dot file")

	return cmd
}
