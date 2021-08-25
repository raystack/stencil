package cmd

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/odpf/stencil/server/api/v1/genproto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var namespace, name, version, filePath string
var latest bool

func init() {
	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload filedescriptorset file",
		Run:   upload,
		Args:  cobra.NoArgs,
	}
	uploadCmd.Flags().StringVarP(&namespace, "namespace", "g", "", "provide namespace/group or entity name")
	uploadCmd.MarkFlagRequired("group")
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
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := genproto.NewStencilServiceClient(conn)
	s := &genproto.Snapshot{
		Namespace: namespace,
		Name:      name,
		Version:   version,
		Latest:    latest,
	}
	ur := &genproto.UploadRequest{
		Snapshot: s,
		Data:     fileData,
	}
	res, err := client.Upload(context.Background(), ur)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res)
}
