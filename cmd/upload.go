package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	stencilv1beta1 "github.com/odpf/stencil/server/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UploadCmd creates a new cobra command for upload
func UploadCmd() *cobra.Command {

	var host, filePath string
	var req stencilv1beta1.CreateSchemaRequest
	var format, compatibility string

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload filedescriptorset file",
		Args:  cobra.NoArgs,
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
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
			client := stencilv1beta1.NewStencilServiceClient(conn)
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[compatibility])
			req.Format = stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[format])
			_, err = client.CreateSchema(context.Background(), &req)
			if err != nil {
				errStatus := status.Convert(err)
				return errors.New(errStatus.Message())
			}
			fmt.Println("success")
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")
	cmd.Flags().StringVar(&req.NamespaceId, "namespace", "", "provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")
	cmd.Flags().StringVar(&req.SchemaId, "name", "", "provide proto repo name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&filePath, "file", "", "provide path to fully contained file descriptor set file")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringVar(&format, "format", "", "schema format. Valid values are FORMAT_PROTOBUF,FORMAT_AVRO,FORMAT_JSON")
	cmd.MarkFlagRequired("format")
	cmd.Flags().StringVar(&compatibility, "compatibility", "COMPATIBILITY_FULL", "schema compatibility. Valid values are COMPATIBILITY_FULL")
	cmd.MarkFlagRequired("format")
	return cmd
}
