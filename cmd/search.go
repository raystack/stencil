package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/printer"
	stencilv1beta1 "github.com/odpf/stencil/server/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func SearchCmd() *cobra.Command {
	var host, namespaceID, schemaID string
	var versionID int32
	var history bool
	var req stencilv1beta1.SearchRequest

	cmd := &cobra.Command{
		Use:     "search",
		Aliases: []string{"search"},
		Short:   "Search",
		Long:    "Search your queries on schemas",
		Example: heredoc.Doc(`
			$ stencil search <query> --namespace=<namespace> --schema=<schema> --version=<version> --history=<history>
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

			query := args[0]
			req.Query = query
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

			client := stencilv1beta1.NewStencilServiceClient(conn)
			res, err := client.Search(context.Background(), &req)
			if err != nil {
				return err
			}

			hits := res.GetHits()
			meta := res.GetMeta()

			report := [][]string{}
			s.Stop()

			if len(hits) == 0 {
				fmt.Printf("No results found")
				return nil
			}

			fmt.Printf(" \nShowing %d result(s) \n", len(hits))

			report = append(report, []string{"INDEX", "NAMESPACE", "SCHEMA", "FIELDS", "TYPES", "VERSION"})
			for i, h := range hits {
				report = append(report, []string{
					strconv.Itoa(i + 1),
					h.GetNamespaceId(),
					h.GetSchemaId(),
					strings.Join(h.GetFields(), ","),
					strings.Join(h.GetTypes(), ","),
					strconv.Itoa(int(h.GetVersionId())),
				})
			}
			printer.Table(os.Stdout, report)

			report = [][]string{}
			report = append(report, []string{"\nTOTAL"})
			report = append(report, []string{strconv.Itoa(int(meta.GetTotal()))})
			printer.Table(os.Stdout, report)

			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")
	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "parent namespace ID")
	cmd.MarkFlagRequired("namespace")
	cmd.Flags().StringVarP(&schemaID, "schema", "s", "", "related schema ID")
	cmd.MarkFlagRequired("schema")
	cmd.Flags().Int32VarP(&versionID, "version", "v", 0, "version of the schema")
	cmd.Flags().BoolVarP(&history, "history", "h", false, "set this to enable history")

	return cmd
}
