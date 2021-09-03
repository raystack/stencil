package cmd

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/odpf/stencil/server/api/v1/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var host, namespace, name, version, filePath string
var latest bool

func init() {
	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload filedescriptorset file",
		Run:   upload,
		Args:  cobra.NoArgs,
	}
	uploadCmd.Flags().StringVarP(&host, "host", "h", "", "stencil host address eg: localhost:8000")
	uploadCmd.MarkFlagRequired("host")
	uploadCmd.Flags().StringVarP(&namespace, "namespace", "g", "", "provide namespace/group or entity name")
	uploadCmd.MarkFlagRequired("namespace")
	uploadCmd.Flags().StringVarP(&name, "name", "n", "", "provide proto repo name")
	uploadCmd.MarkFlagRequired("name")
	uploadCmd.Flags().StringVarP(&version, "version", "v", "", "provide semantic version compatible value")
	uploadCmd.MarkFlagRequired("version")
	uploadCmd.Flags().BoolVarP(&latest, "latest", "l", false, "mark as latest version")
	uploadCmd.Flags().StringVarP(&filePath, "file", "f", "", "provide path to fully contained file descriptor set file")
	uploadCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(uploadCmd)
}

func upload(cmd *cobra.Command, args []string) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Unable to read provided file", err)
	}
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := pb.NewStencilServiceClient(conn)
	ur := &pb.UploadDescriptorRequest{
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
}
