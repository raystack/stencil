package cmd

import (
	"context"
	"errors"
	"fmt"

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
	}
	cmd.PersistentFlags().String("host", "", "stencil host address eg: localhost:8000")
	cmd.MarkPersistentFlagRequired("host")
	cmd.AddCommand(listCmd())
	cmd.AddCommand(promoteCmd())
	return cmd
}

func listCmd() *cobra.Command {
	var req stencilv1.ListSnapshotsRequest
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list snapshots with optional filters",
		Args:  cobra.NoArgs,
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
