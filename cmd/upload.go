package cmd

import (
	"context"
	"io/ioutil"
	"log"

	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// UploadCmd creates a new cobra command for upload
func UploadCmd() *cobra.Command {

	var host, namespace, name, version, filePath string
	var latest bool

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload filedescriptorset file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fileData, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Fatalln("Unable to read provided file", err)
			}
			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				log.Fatalln(err)
			}
			defer conn.Close()
			client := stencilv1.NewStencilServiceClient(conn)
			ur := &stencilv1.UploadDescriptorRequest{
				Namespace: namespace,
				Name:      name,
				Version:   version,
				Latest:    latest,
				Data:      fileData,
			}
			res, err := client.UploadDescriptor(context.Background(), ur)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(res)
			return nil
		},
	}

	cmd.Flags().StringVarP(&host, "host", "h", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")
	cmd.Flags().StringVarP(&namespace, "namespace", "g", "", "provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")
	cmd.Flags().StringVarP(&name, "name", "n", "", "provide proto repo name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&version, "version", "v", "", "provide semantic version compatible value")
	cmd.MarkFlagRequired("version")
	cmd.Flags().BoolVarP(&latest, "latest", "l", false, "mark as latest version")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "provide path to fully contained file descriptor set file")
	cmd.MarkFlagRequired("file")
	return cmd
}
