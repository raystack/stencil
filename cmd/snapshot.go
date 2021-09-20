package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// Snapshot creates a new cobra command to manage snapshot
func Snapshot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "list, update snapshot details",
		Annotations: map[string]string{
			"group:core": "true",
		},
	}
	cmd.PersistentFlags().String("host", "", "stencil host address eg: localhost:8000")
	cmd.MarkPersistentFlagRequired("host")
	cmd.AddCommand(listCmd())
	cmd.AddCommand(promoteCmd())
	cmd.AddCommand(printCmd())
	return cmd
}

func listCmd() *cobra.Command {
	var req stencilv1.ListSnapshotsRequest
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list snapshots with optional filters",
		Args:  cobra.NoArgs,
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			host, _ := cmd.Flags().GetString("host")
			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			client := stencilv1.NewStencilServiceClient(conn)
			res, err := client.ListSnapshots(context.Background(), &req)
			if err != nil {
				errStatus := status.Convert(err)
				return errors.New(errStatus.Message())
			}
			data, _ := protojson.MarshalOptions{EmitUnpopulated: true, Multiline: true, Indent: "  "}.Marshal(res)
			fmt.Println(string(data))
			return nil
		},
	}
	cmd.Flags().StringVar(&req.Namespace, "namespace", "", "provide namespace/group or entity name")
	cmd.Flags().StringVar(&req.Name, "name", "", "provide proto repo name")
	cmd.Flags().StringVar(&req.Version, "version", "", "provide semantic version compatible value")
	cmd.Flags().BoolVar(&req.Latest, "latest", false, "mark as latest version")
	return cmd
}

func promoteCmd() *cobra.Command {
	var req stencilv1.PromoteSnapshotRequest
	cmd := &cobra.Command{
		Use:   "promote",
		Short: "promote specified snapshot to latest",
		Args:  cobra.NoArgs,
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			host, _ := cmd.Flags().GetString("host")
			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			client := stencilv1.NewStencilServiceClient(conn)
			res, err := client.PromoteSnapshot(context.Background(), &req)
			if err != nil {
				errStatus := status.Convert(err)
				return errors.New(errStatus.Message())
			}
			data, _ := protojson.MarshalOptions{EmitUnpopulated: true, Multiline: true, Indent: "  "}.Marshal(res)
			fmt.Println(string(data))
			return nil
		},
	}
	cmd.Flags().Int64Var(&req.Id, "id", 0, "snapshot id")
	cmd.MarkFlagRequired("id")
	return cmd
}

// printCmd creates a new cobra command for print
func printCmd() *cobra.Command {
	var (
		req              stencilv1.DownloadDescriptorRequest
		pathDir          string
		filterPathPrefix string
	)
	cmd := &cobra.Command{
		Use:   "print",
		Short: "prints snapshot details into .proto files",
		Args:  cobra.NoArgs,
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			host, _ := cmd.Flags().GetString("host")
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

			fds := &descriptorpb.FileDescriptorSet{}
			if err := proto.Unmarshal(res.GetData(), fds); err != nil {
				return fmt.Errorf("descriptor set file is not valid. %w", err)
			}
			fdsMap, err := desc.CreateFileDescriptorsFromSet(fds)
			if err != nil {
				return err
			}

			var filteredFds []*desc.FileDescriptor
			for fdName, fd := range fdsMap {
				if filterPathPrefix != "" && !strings.HasPrefix(fdName, filterPathPrefix) {
					continue
				}
				filteredFds = append(filteredFds, fd)
			}

			protoPrinter := &protoprint.Printer{}
			if pathDir == "" {
				for _, fd := range filteredFds {
					protoAsString, err := protoPrinter.PrintProtoToString(fd)
					if err != nil {
						return err
					}
					fmt.Printf("\n// ----\n// %s\n// ----\n%s", fd.GetName(), protoAsString)
				}
			} else {
				if err := protoPrinter.PrintProtosToFileSystem(filteredFds, pathDir); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&req.Namespace, "namespace", "", "provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")
	cmd.Flags().StringVar(&req.Name, "name", "", "provide proto repo name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&req.Version, "version", "", "provide semantic version compatible value")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringVar(&pathDir, "output", "", "the directory path to write the descriptor files, default is to print on stdout")
	cmd.Flags().StringVar(&filterPathPrefix, "filter-path", "", "filter protocol buffer files by path prefix, e.g., --filter-path=google/protobuf")
	return cmd
}
