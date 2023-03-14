package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/alecthomas/chroma/quick"
	stencilv1beta1 "github.com/goto/stencil/proto/v1beta1"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/odpf/salt/printer"
	"github.com/odpf/salt/term"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func printSchemaCmd(cdk *CDK) *cobra.Command {
	var filter, namespaceID string
	var version int32

	cmd := &cobra.Command{
		Use:     "view <id>",
		Short:   "Print snapshot of a schema",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"print"},
		Example: heredoc.Doc(`
			$ stencil schema view booking -n goto
			$ stencil schema view booking -n goto -v 2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()
			client, cancel, err := createClient(cmd, cdk)
			if err != nil {
				return err
			}
			defer cancel()

			data, meta, err := fetchSchemaAndMeta(client, version, namespaceID, args[0])
			if err != nil {
				return err
			}
			spinner.Stop()

			format := stencilv1beta1.Schema_Format_name[int32(meta.GetFormat())]
			switch format {
			case "FORMAT_AVRO":
				if err := printSchema(data); err != nil {
					return err
				}
			case "FORMAT_JSON":
				if err := printSchema(data); err != nil {
					return err
				}
			case "FORMAT_PROTOBUF":
				printProtoSchema(data, filter)
			default:
				fmt.Printf("%s Unknown schema format: %s\n", term.Red(term.FailureIcon()), format)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "Provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().Int32VarP(&version, "version", "v", 0, "Provide version number")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter schema files by path prefix, e.g., --filter=google/protobuf")

	return cmd
}

func printSchema(data []byte) error {
	page := term.NewPager()
	page.Start()
	defer page.Stop()

	err := quick.Highlight(page.Out, string(data), "JSON", "terminal16m", "solarized-light")
	if err != nil {
		page.Out.Write(data)
	}
	return nil
}

func printProtoSchema(data []byte, filter string) error {
	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(data, fds); err != nil {
		return fmt.Errorf("descriptor set file is not valid. %w", err)
	}
	fdsMap, err := desc.CreateFileDescriptorsFromSet(fds)
	if err != nil {
		return err
	}
	var filteredFds []*desc.FileDescriptor
	for fdName, fd := range fdsMap {
		if filter != "" && !strings.HasPrefix(fdName, filter) {
			continue
		}
		filteredFds = append(filteredFds, fd)
	}

	protoPrinter := &protoprint.Printer{}

	var schema string

	for _, fd := range filteredFds {
		protoAsString, err := protoPrinter.PrintProtoToString(fd)
		if err != nil {
			return err
		}
		schema = schema + fmt.Sprintf("\n//Schema file:: %s\n\n%s", fd.GetName(), protoAsString)
	}

	page := term.NewPager()
	page.Start()
	defer page.Stop()

	err = quick.Highlight(page.Out, schema, "Protocol Buffer", "terminal16m", "solarized-light")
	if err != nil {
		fmt.Fprint(page.Out, schema)
	}
	return nil
}
