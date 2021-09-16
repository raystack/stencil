package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UploadCmd creates a new cobra command for upload
func UploadCmd() *cobra.Command {

	var host, filePath string
	var req stencilv1.UploadDescriptorRequest
	var skipRules []string

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
			req.Checks = &stencilv1.Checks{
				Except: toRules(skipRules),
			}
			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			client := stencilv1.NewStencilServiceClient(conn)
			_, err = client.UploadDescriptor(context.Background(), &req)
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
	cmd.Flags().StringVar(&req.Namespace, "namespace", "", "provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")
	cmd.Flags().StringVar(&req.Name, "name", "", "provide proto repo name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&req.Version, "version", "", "provide semantic version compatible value")
	cmd.MarkFlagRequired("version")
	cmd.Flags().BoolVar(&req.Latest, "latest", false, "mark as latest version")
	cmd.Flags().StringVar(&filePath, "file", "", "provide path to fully contained file descriptor set file")
	cmd.MarkFlagRequired("file")
	cmd.Flags().BoolVar(&req.Dryrun, "dryrun", false, "enable dryrun flag")
	cmd.Flags().StringArrayVar(&skipRules, "skiprules", []string{}, "list of rules to skip. Invalid rules ignored Eg: FILE_NO_BREAKING_CHANGE")
	return cmd
}

func toRules(stringRules []string) []stencilv1.Rule {
	var rules []stencilv1.Rule
	for _, rule := range stringRules {
		if val, ok := stencilv1.Rule_value[rule]; ok {
			rules = append(rules, stencilv1.Rule(val))
		}
	}
	return rules
}
