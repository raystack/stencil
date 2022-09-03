package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/dustin/go-humanize"
	"github.com/odpf/salt/printer"
	"github.com/odpf/salt/term"
	"github.com/odpf/stencil/pkg/prompt"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			$ stencil namespace view odpf
			$ stencil namespace delete odpf
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
	var host, id, desc, format, comp string
	var req stencilv1beta1.CreateNamespaceRequest

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a namespace",
		Args:  cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			$ stencil namespace create 
			$ stencil namespace create -n=odpf -f=FORMAT_PROTOBUF --c=COMPATIBILITY_BACKWARD --d="Event schemas"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompter := prompt.New()
			if id == "" {
				id, _ = prompter.Input("What is the namespace id?", "")
			}

			if desc == "" {
				desc, _ = prompter.Input("Provide a description?", "")
			}

			if format == "" {
				formats := []string{"FORMAT_JSON", "FORMAT_PROTOBUF", "FORMAT_AVRO"}
				formatAnswer, _ := prompter.Select("Select a default schema format for this namespace:", formats[0], formats)
				format = formats[formatAnswer]
			}

			if comp == "" {
				comps := []string{"COMPATIBILITY_BACKWARD", "COMPATIBILITY_FORWARD", "COMPATIBILITY_FULL"}
				formatAnswer, _ := prompter.Select("Select a default compatibility for this namespace:", comps[0], comps)
				fmt.Println()
				comp = comps[formatAnswer]
			}

			req.Id = id
			req.Format = stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[format])
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp])
			req.Description = desc

			spinner := printer.Spin("")
			defer spinner.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			client := stencilv1beta1.NewStencilServiceClient(conn)
			res, err := client.CreateNamespace(context.Background(), &req)
			spinner.Stop()

			if err != nil {
				errStatus, _ := status.FromError(err)
				if codes.AlreadyExists == errStatus.Code() {
					fmt.Printf("\n%s Namespace with id '%s' already exist.\n", term.FailureIcon(), id)
					return nil
				}
				return err
			}

			namespace := res.GetNamespace()
			fmt.Printf("\n%s Created namespace with id '%s'.\n", term.Green(term.SuccessIcon()), namespace.GetId())
			return nil
		},
	}

	cmd.Flags().StringVarP(&id, "id", "n", "", "Supply an id. Will prompt otherwise")
	cmd.Flags().StringVarP(&format, "format", "f", "", "Default schema format for schemas in this namespace")
	cmd.Flags().StringVarP(&comp, "comp", "c", "", "Default schema compatibility for schemas in this namespace")
	cmd.Flags().StringVarP(&desc, "desc", "d", "", "Supply a description. Will prompt otherwise")

	cmd.Flags().StringVar(&host, "host", "", "Stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

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
			res, err := client.UpdateNamespace(context.Background(), &req)
			spinner.Stop()
			if err != nil {
				errStatus, _ := status.FromError(err)
				if codes.NotFound == errStatus.Code() {
					fmt.Printf("%s Namespace with id '%s' does not exist.\n", term.FailureIcon(), id)
					return nil
				}
				return err
			}

			namespace := res.Namespace

			fmt.Printf("%s Updated namespace with id '%s'.\n", term.Green(term.SuccessIcon()), namespace.GetId())
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
		Use:   "view <id>",
		Short: "View a namespace",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace view <id>
		`),
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
			spinner.Stop()

			if err != nil {
				errStatus, _ := status.FromError(err)
				if codes.NotFound == errStatus.Code() {
					fmt.Printf("%s Namespace with id '%s' does not exist.\n", term.FailureIcon(), id)
					return nil
				}
				return err
			}

			namespace := res.GetNamespace()

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
		Use:   "delete <id>",
		Short: "Delete a namespace",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace delete <id>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			prompter := prompt.New()
			confirm, _ := prompter.Input(fmt.Sprintf("You're going to delete namespace `%s`. To confirm, type the namespace id:", id), "")
			if id != confirm {
				fmt.Printf("\n%s Namespace id '%s' did not match.\n", term.WarningIcon(), confirm)
				return nil
			}

			spinner := printer.Spin("")
			defer spinner.Stop()

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			req.Id = id

			client := stencilv1beta1.NewStencilServiceClient(conn)
			_, err = client.DeleteNamespace(context.Background(), &req)

			spinner.Stop()

			if err != nil {
				errStatus, _ := status.FromError(err)
				if codes.NotFound == errStatus.Code() {
					fmt.Printf("\n%s Namespace with id '%s' does not exist.\n", term.FailureIcon(), id)
					return nil
				}
				return err
			}

			fmt.Printf("\n%s Deleted namespace with id '%s'.\n", term.Red(term.SuccessIcon()), id)

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
