package cmd

import (
	"context"
	"os"

	stencilv1beta1 "github.com/odpf/stencil/server/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// DownloadCmd creates a new cobra command for download descriptor
func DownloadCmd() *cobra.Command {
	var host, filePath string
	var req stencilv1beta1.GetSchemaRequest

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download filedescriptorset file",
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
			client := stencilv1beta1.NewStencilServiceClient(conn)
			res, err := client.GetSchema(context.Background(), &req)
			if err != nil {
				return err
			}
			err = os.WriteFile(filePath, res.Data, 0666)
			return err
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
	cmd.Flags().StringVar(&filePath, "output", "", "write to file")
	return cmd
}
