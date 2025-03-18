package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/dustin/go-humanize"
	"github.com/raystack/salt/cli/printer"
	"github.com/raystack/salt/cli/prompter"
	stencilv1beta1 "github.com/raystack/stencil/proto/raystack/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NamespaceCmd(cdk *CDK) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "namespace",
		Aliases: []string{"namespaces"},
		Short:   "Manage namespaces",
		Long:    "Work with namespaces.",
		Example: heredoc.Doc(`
			$ stencil namespace list
			$ stencil namespace create -n raystack
			$ stencil namespace view raystack
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(listNamespaceCmd(cdk))
	cmd.AddCommand(createNamespaceCmd(cdk))
	cmd.AddCommand(viewNamespaceCmd(cdk))
	cmd.AddCommand(editNamespaceCmd(cdk))
	cmd.AddCommand(deleteNamespaceCmd(cdk))

	return cmd
}

func listNamespaceCmd(cdk *CDK) *cobra.Command {
	var req stencilv1beta1.ListNamespacesRequest

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all namespaces",
		Long:  "List and filter namespaces.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()
			client, cancel, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.ListNamespaces(context.Background(), &req)
			if err != nil {
				return err
			}

			namespaces := res.GetNamespaces()
			spinner.Stop()

			if len(namespaces) == 0 {
				fmt.Println("No namespace found")
				return nil
			}

			fmt.Printf("\nShowing %[1]d of %[1]d namespaces \n \n", len(namespaces))
			report := [][]string{}
			index := 1
			report = append(report, []string{
				printer.Bold("INDEX"),
				printer.Bold("NAMESPACE"),
				printer.Bold("FORMAT"),
				printer.Bold("COMPATIBILITY"),
			})
			for _, n := range namespaces {
				report = append(report,
					[]string{printer.Greenf("#%d", index),
						n.Id,
						dict[n.GetFormat().String()],
						dict[n.GetCompatibility().String()],
					})
				index++
			}
			printer.Table(os.Stdout, report)
			return nil
		},
	}

	return cmd
}

func createNamespaceCmd(cdk *CDK) *cobra.Command {
	var id, desc, format, comp string
	var req stencilv1beta1.CreateNamespaceRequest

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a namespace",
		Args:  cobra.ExactArgs(0),
		Example: heredoc.Doc(`
			$ stencil namespace create 
			$ stencil namespace create -n=raystack -f=FORMAT_PROTOBUF -c=COMPATIBILITY_BACKWARD -d="Event schemas"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompter := prompter.New()
			if id == "" {
				id, _ = prompter.Input("What is the namespace id?", "")
			}

			if desc == "" {
				desc, _ = prompter.Input("Provide a description?", "")
			}

			if format == "" {
				formatAnswer, _ := prompter.Select("Select a default schema format for this namespace:", formats[0], formats)
				format = formats[formatAnswer]
			}

			if comp == "" {
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

			client, cancel, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.CreateNamespace(context.Background(), &req)
			spinner.Stop()

			if err != nil {
				errStatus, _ := status.FromError(err)
				if codes.AlreadyExists == errStatus.Code() {
					fmt.Printf("\n%s Namespace with id '%s' already exist.\n", printer.Icon("failure"), id)
					return nil
				}
				return err
			}

			namespace := res.GetNamespace()
			fmt.Printf("\n%s Created namespace with id %s.\n", printer.Green(printer.Icon("success")), printer.Bold(printer.Blue(namespace.GetId())))
			return nil
		},
	}

	cmd.Flags().StringVarP(&id, "id", "n", "", "Supply an id. Will prompt otherwise")
	cmd.Flags().StringVarP(&format, "format", "f", "", "Default schema format for schemas in this namespace")
	cmd.Flags().StringVarP(&comp, "comp", "c", "", "Default schema compatibility for schemas in this namespace")
	cmd.Flags().StringVarP(&desc, "desc", "d", "", "Supply a description. Will prompt otherwise")

	return cmd
}

func editNamespaceCmd(cdk *CDK) *cobra.Command {
	var format, comp string
	var desc string
	var req stencilv1beta1.UpdateNamespaceRequest

	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit a namespace",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace edit raystack -f FORMAT_JSON -c COMPATIBILITY_BACKWARD -d "Hello message"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			client, cancel, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}
			defer cancel()

			id := args[0]

			req.Id = id
			req.Format = stencilv1beta1.Schema_Format(stencilv1beta1.Schema_Format_value[format])
			req.Compatibility = stencilv1beta1.Schema_Compatibility(stencilv1beta1.Schema_Compatibility_value[comp])
			req.Description = desc

			res, err := client.UpdateNamespace(context.Background(), &req)
			spinner.Stop()

			if err != nil {
				errStatus, _ := status.FromError(err)
				if codes.NotFound == errStatus.Code() {
					fmt.Printf("%s Namespace with id '%s' does not exist.\n", printer.Icon("failure"), id)
					return nil
				}
				return err
			}

			namespace := res.Namespace

			fmt.Printf("%s Updated namespace with id %s.\n", printer.Green(printer.Icon("success")), printer.Bold(printer.Blue(namespace.GetId())))
			return nil
		},
	}

	// TODO(Ravi) : Edit should not require all flags
	cmd.Flags().StringVarP(&format, "format", "f", "", "schema format")
	cmd.MarkFlagRequired("format")

	cmd.Flags().StringVarP(&comp, "comp", "c", "", "schema compatibility")
	cmd.MarkFlagRequired("comp")

	cmd.Flags().StringVarP(&desc, "desc", "d", "", "description")
	cmd.MarkFlagRequired("desc")

	return cmd
}

func viewNamespaceCmd(cdk *CDK) *cobra.Command {
	var req stencilv1beta1.GetNamespaceRequest

	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "View a namespace",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace view raystack
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			client, cancel, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}
			defer cancel()

			id := args[0]
			req.Id = id

			res, err := client.GetNamespace(context.Background(), &req)
			spinner.Stop()

			if err != nil {
				errStatus, _ := status.FromError(err)
				if codes.NotFound == errStatus.Code() {
					fmt.Printf("%s Namespace with id %s does not exist.\n", printer.Icon("failure"), printer.Bold(printer.Blue(id)))
					return nil
				}
				return err
			}

			namespace := res.GetNamespace()

			printNamespace(namespace)

			return nil
		},
	}

	return cmd
}

func deleteNamespaceCmd(cdk *CDK) *cobra.Command {
	var req stencilv1beta1.DeleteNamespaceRequest

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a namespace",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil namespace delete raystack
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			prompter := prompter.New()
			confirm, _ := prompter.Input(fmt.Sprintf("Deleting namespace `%s`. To confirm, type the namespace id:", id), "")
			if id != confirm {
				fmt.Printf("\n%s Namespace id '%s' did not match.\n", printer.Icon("warning"), confirm)
				return nil
			}

			spinner := printer.Spin("")
			defer spinner.Stop()

			client, cancel, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}
			defer cancel()

			req.Id = id

			_, err = client.DeleteNamespace(context.Background(), &req)

			spinner.Stop()

			if err != nil {
				errStatus, _ := status.FromError(err)
				if codes.NotFound == errStatus.Code() {
					fmt.Printf("\n%s Namespace with id '%s' does not exist.\n", printer.Icon("failure"), id)
					return nil
				}
				return err
			}

			fmt.Printf("\n%s Deleted namespace with id %s.\n", printer.Red(printer.Icon("success")), printer.Bold(printer.Blue(id)))

			return nil
		},
	}

	return cmd
}

func printNamespace(namespace *stencilv1beta1.Namespace) {
	desc := namespace.GetDescription()
	if desc == "" {
		desc = "No description provided"
	}

	fmt.Printf("\n%s\n", printer.Blue(namespace.GetId()))
	fmt.Printf("\n%s.\n\n", printer.Grey(desc))
	fmt.Printf("%s \t %s \n", printer.Grey("Format:"), namespace.GetFormat().String())
	fmt.Printf("%s \t %s \n", printer.Grey("Compatibility:"), namespace.GetCompatibility().String())
	fmt.Printf("\n%s %s, ", printer.Grey("Created"), humanize.Time(namespace.GetCreatedAt().AsTime()))
	fmt.Printf("%s %s \n\n", printer.Grey("last updated"), humanize.Time(namespace.GetUpdatedAt().AsTime()))
}
