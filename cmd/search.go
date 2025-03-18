package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/raystack/salt/cli/printer"
	stencilv1beta1 "github.com/raystack/stencil/proto/raystack/stencil/v1beta1"
	"github.com/spf13/cobra"
)

func SearchCmd(cdk *CDK) *cobra.Command {
	var namespaceID, schemaID string
	var versionID int32
	var history bool
	var req stencilv1beta1.SearchRequest

	cmd := &cobra.Command{
		Use:     "search <query>",
		Aliases: []string{"search"},
		Short:   "Search schemas",
		Long:    "Search your queries on schemas",
		Args:    cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil search email
			$ stencil search email -s human
			$ stencil search name -n raystack -s person -v 2
			$ stencil search address -n raystack -s person -h true
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		s := printer.Spin("")
		defer s.Stop()

		client, cancel, err := createClient(cmd, cdk)
		if err != nil {
			return err
		}
		defer cancel()

		query := args[0]
		req.Query = query

		if len(schemaID) > 0 && len(namespaceID) == 0 {
			s.Stop()
			fmt.Println("Namespace ID not specified for", schemaID)
			return nil
		}
		req.NamespaceId = namespaceID
		req.SchemaId = schemaID

		if versionID != 0 {
			req.Version = &stencilv1beta1.SearchRequest_VersionId{
				VersionId: versionID,
			}
		} else if history {
			req.Version = &stencilv1beta1.SearchRequest_History{
				History: history,
			}
		}

		res, err := client.Search(context.Background(), &req)
		if err != nil {
			return err
		}

		hits := res.GetHits()

		report := [][]string{}
		s.Stop()

		if len(hits) == 0 {
			fmt.Println("No results found")
			return nil
		}

		var total = 0
		report = append(report, []string{
			printer.Bold("FIELD"),
			printer.Bold("TYPE"),
			printer.Bold("SCHEMA"),
			printer.Bold("VERSION"),
			printer.Bold("NAMESPACE"),
		})
		for _, h := range hits {
			fields := h.GetFields()
			for _, field := range fields {
				report = append(report, []string{
					field[strings.LastIndex(field, ".")+1:],
					field[:strings.LastIndex(field, ".")],
					h.GetSchemaId(),
					strconv.Itoa(int(h.GetVersionId())),
					h.GetNamespaceId(),
				})
				total++
			}
		}
		fmt.Printf(" \nFound %d results across %d schema(s)/version(s) \n\n", total, len(hits))
		printer.Table(os.Stdout, report)
		return nil
	}

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "parent namespace ID")
	cmd.Flags().StringVarP(&schemaID, "schema", "s", "", "related schema ID")
	cmd.Flags().Int32VarP(&versionID, "version", "v", 0, "version of the schema")
	cmd.Flags().BoolVarP(&history, "history", "h", false, "set this to enable history")

	return cmd
}
