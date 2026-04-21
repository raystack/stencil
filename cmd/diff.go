package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"connectrpc.com/connect"
	"github.com/MakeNowJust/heredoc"
	"github.com/raystack/salt/cli/printer"
	stencilv1beta1 "github.com/raystack/stencil/gen/raystack/stencil/v1beta1"
	stencilv1beta1connect "github.com/raystack/stencil/gen/raystack/stencil/v1beta1/stencilv1beta1connect"
	"github.com/spf13/cobra"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func diffSchemaCmd(cdk *CDK) *cobra.Command {
	var fullname string
	var namespace string
	var earlierVersion int32
	var laterVersion int32

	var schemaFetcher = func(req *stencilv1beta1.GetSchemaRequest, client stencilv1beta1connect.StencilServiceClient) ([]byte, error) {
		res, err := client.GetSchema(context.Background(), connect.NewRequest(req))
		if err != nil {
			return nil, err
		}
		return res.Msg.GetData(), nil
	}
	var protoSchemaFetcher = func(req *stencilv1beta1.GetSchemaRequest, client stencilv1beta1connect.StencilServiceClient) ([]byte, error) {
		if fullname == "" {
			return nil, fmt.Errorf("fullname flag is mandator for FORMAT_PROTO")
		}
		res, err := client.GetSchema(context.Background(), connect.NewRequest(req))
		if err != nil {
			return nil, err
		}
		fds := &descriptorpb.FileDescriptorSet{}
		if err := proto.Unmarshal(res.Msg.GetData(), fds); err != nil {
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
			$ stencil schema diff booking -n=raystack --later-version=2 --earlier-version=1 --fullname=<fullname>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			schemaID := args[0]

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

			client, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}

			metaRes, err := client.GetSchemaMetadata(context.Background(), connect.NewRequest(&stencilv1beta1.GetSchemaMetadataRequest{
				NamespaceId: namespace,
				SchemaId:    schemaID,
			}))
			if err != nil {
				return err
			}

			meta := metaRes.Msg
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
			if err := json.Unmarshal(eJson, &placeholder); err != nil {
				return err
			}
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

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Parent namespace ID")
	_ = cmd.MarkFlagRequired("namespace")
	cmd.Flags().Int32Var(&earlierVersion, "earlier-version", 0, "Earlier version of the schema")
	_ = cmd.MarkFlagRequired("earlier-version")
	cmd.Flags().Int32Var(&laterVersion, "later-version", 0, "Later version of the schema")
	_ = cmd.MarkFlagRequired("later-version")
	cmd.Flags().StringVar(&fullname, "fullname", "", "Only applicable for FORMAT_PROTO. fullname of proto schema eg: raystack.common.v1.Version")
	return cmd
}
