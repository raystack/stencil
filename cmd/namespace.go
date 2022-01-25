package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	stencilv1beta1 "github.com/odpf/stencil/server/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NamespaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "namespace",
		Aliases: []string{"namespace"},
		Short:   "Manage namespace",
		Long: heredoc.Doc(`
			Work with namespaces.
		`),
		Example: heredoc.Doc(`
			$ stencil namespace list
			$ stencil namespace create
			$ stencil namespace get
			$ stencil namespace update
			$ stencil namespace delete
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
	}

	cmd.AddCommand(listNamespaceCmd())
	cmd.AddCommand(createNamespaceCmd())
	cmd.AddCommand(getNamespaceCmd())
	cmd.AddCommand(updateNamespaceCmd())
	cmd.AddCommand(deleteNamespaceCmd())

	return cmd
}

func listNamespaceCmd() *cobra.Command {
	var host string
	var req stencilv1beta1.ListNamespacesRequest

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all namespaces",
		Args:  cobra.NoArgs,
		Example: heredoc.Doc(`
			$ stencil namespace list
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := printer.Spin("")
			defer s.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			client := stencilv1beta1.NewStencilServiceClient(conn)
			res, err := client.ListNamespaces(context.Background(), &req)
			if err != nil {
				return err
			}

			report := [][]string{}

			namespaces := res.GetNamespaces()

			s.Stop()

			fmt.Printf(" \nShowing %d namespaces \n", len(namespaces))

			report = append(report, []string{"NAMESPACE"})

			for _, n := range namespaces {
				report = append(report, []string{
					n,
				})
			}
			printer.Table(os.Stdout, report)
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	return cmd
}

func createNamespaceCmd() *cobra.Command {
	var host, format, comp string
	var desc string
	var req stencilv1beta1.CreateNamespaceRequest

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create all namespaces",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace create <namespace-id> --format=<schema-format> –-comp=<schema-compatibility> –-desc=<description> 
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			id := args[0]

			req.Id = id
			req.Format = stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[format])
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp])
			req.Description = desc

			client := stencilv1beta1.NewStencilServiceClient(conn)
			res, err := client.CreateNamespace(context.Background(), &req)
			if err != nil {
				return err
			}

			namespace := res.GetNamespace()

			spinner.Stop()

			fmt.Printf("namespace successfully created with id: %s", namespace.GetId())
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&format, "format", "f", "", "schema format")
	cmd.MarkFlagRequired("format")

	cmd.Flags().StringVarP(&comp, "comp", "c", "", "schema compatibility")
	cmd.MarkFlagRequired("comp")

	cmd.Flags().StringVarP(&desc, "desc", "d", "", "description")

	return cmd
}

func updateNamespaceCmd() *cobra.Command {
	var host, format, comp string
	var desc string
	var req stencilv1beta1.UpdateNamespaceRequest

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update all namespaces",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace update <namespace-id> –-format=<schema-format> –-comp=<schema-compatibility> –-desc=<description>
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			id := args[0]

			req.Id = id
			req.Format = stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[format])
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp])
			req.Description = desc

			client := stencilv1beta1.NewStencilServiceClient(conn)
			_, err = client.UpdateNamespace(context.Background(), &req)
			if err != nil {
				return err
			}

			spinner.Stop()

			fmt.Printf("namespace successfully updated")
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&format, "format", "f", "", "schema format")
	cmd.MarkFlagRequired("format")

	cmd.Flags().StringVarP(&comp, "comp", "c", "", "schema compatibility")
	cmd.MarkFlagRequired("comp")

	cmd.Flags().StringVarP(&desc, "desc", "d", "", "description")
	cmd.MarkFlagRequired("desc")

	return cmd
}

func getNamespaceCmd() *cobra.Command {
	var host string
	var req stencilv1beta1.GetNamespaceRequest

	cmd := &cobra.Command{
		Use:   "get",
		Short: "get all namespaces",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace get <namespace-id>
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			id := args[0]

			req.Id = id

			client := stencilv1beta1.NewStencilServiceClient(conn)
			res, err := client.GetNamespace(context.Background(), &req)
			if err != nil {
				return err
			}

			report := [][]string{}

			namespace := res.GetNamespace()

			spinner.Stop()

			report = append(report, []string{"ID", "FORMAT", "COMPATIBILITY", "DESCRIPTION"})
			report = append(report, []string{
				namespace.GetId(),
				namespace.GetFormat().String(),
				namespace.GetCompatibility().String(),
				namespace.GetDescription(),
			})
			printer.Table(os.Stdout, report)
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	return cmd
}

func deleteNamespaceCmd() *cobra.Command {
	var host string
	var req stencilv1beta1.DeleteNamespaceRequest

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete all namespaces",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace delete <namespace-id>
		`),
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			id := args[0]

			req.Id = id

			client := stencilv1beta1.NewStencilServiceClient(conn)
			_, err = client.DeleteNamespace(context.Background(), &req)
			if err != nil {
				return err
			}

			spinner.Stop()

			fmt.Printf("namespace successfully deleted")

			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	return cmd
}
