package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/dustin/go-humanize"
	"github.com/odpf/salt/printer"
	"github.com/odpf/salt/term"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NamespaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "namespace",
		Aliases: []string{"namespace"},
		Short:   "Manage namespaces",
		Long:    "Work with namespaces.",
		Example: heredoc.Doc(`
			$ stencil namespace list
			$ stencil namespace create
			$ stencil namespace view
			$ stencil namespace edit
			$ stencil namespace delete
		`),
		Annotations: map[string]string{
			"group": "core",
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
		Short: "List all namespaces",
		Long:  "List and filter namespaces.",
		Args:  cobra.NoArgs,
		Annotations: map[string]string{
			"group": "core",
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

			fmt.Printf(" \nShowing %[1]d of %[1]d namespaces \n \n", len(namespaces))

			report = append(report, []string{"INDEX", "NAMESPACE", "FORMAT", "COMPATIBILITY", "DESCRIPTION"})
			index := 1

			for _, n := range namespaces {
				report = append(report, []string{term.Greenf("#%02d", index), n, "-", "-", "-"})
				index++
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
		Use:   "create <name>",
		Short: "Create a namespace",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace create odpf -f=FORMAT_PROTOBUF --c=COMPATIBILITY_BACKWARD --d="Event schemas"
		`),
		Annotations: map[string]string{
			"group": "core",
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

			fmt.Printf("Namespace successfully created with id: %s", namespace.GetId())
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
		Use:   "edit <id>",
		Short: "Edit a namespace",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace edit <id> --format=<schema-format> --comp=<schema-compatibility> --desc=<description>
		`),
		Annotations: map[string]string{
			"group": "core",
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

			fmt.Println(term.SuccessIcon(), term.Green("Namespace successfully updated"))
			// TODO(Ravi): Print details of updated namespace
			return nil
		},
	}

	// TODO(Ravi) : Edit should not require all flags
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
		Use:   "view <name>",
		Short: "View a namespace",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace view <id>
		`),
		Annotations: map[string]string{
			"group": "core",
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

			namespace := res.GetNamespace()

			spinner.Stop()

			printNamespace(namespace)

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
		Short: "Delete a namespace",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace delete <id>
		`),
		Annotations: map[string]string{
			"group": "core",
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

			fmt.Printf("Namespace successfully deleted")

			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	return cmd
}

func printNamespace(namespace *stencilv1beta1.Namespace) {
	fmt.Printf("%s \t\t %s \n", term.Bold("Name:"), namespace.GetId())
	fmt.Printf("%s \t %s \n", term.Bold("Format:"), namespace.GetFormat().String())
	fmt.Printf("%s \t %s \n", term.Bold("Compatibility:"), namespace.GetCompatibility().String())
	fmt.Printf("%s \t %s \n\n", term.Bold("Description:"), namespace.GetDescription())
	fmt.Printf("%s %s, ", term.Grey("Created"), humanize.Time(namespace.GetCreatedAt().AsTime()))
	fmt.Printf("%s %s \n\n", term.Grey("last updated"), humanize.Time(namespace.GetCreatedAt().AsTime()))
}
