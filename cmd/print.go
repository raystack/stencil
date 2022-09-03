package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/alecthomas/chroma/quick"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/odpf/salt/printer"
	"github.com/odpf/salt/term"
	stencilv1beta1 "github.com/odpf/stencil/proto/odpf/stencil/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func printCmd() *cobra.Command {
	var output, filterPathPrefix, host, namespaceID string
	var version int32

	cmd := &cobra.Command{
		Use:   "print <id>",
		Short: "Print a given schema snapshot",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ stencil schema print events -n odpf
			$ stencil schema print events -n odpf -v 2 -o ./schema
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()
			client, cancel, err := createClient(cmd)
			if err != nil {
				return err
			}
			defer cancel()

			schemaID := args[0]

			data, meta, err := fetchSchemaAndMetadata(client, version, namespaceID, schemaID)
			if err != nil {
				return err
			}
			spinner.Stop()

			page := term.New()
			page.Start()
			defer page.Stop()

			format := stencilv1beta1.Schema_Format_name[int32(meta.GetFormat())]

			switch format {
			case "FORMAT_AVRO":
				if err := printSchema(page.Out, data, output); err != nil {
					return err
				}
			case "FORMAT_JSON":
				if err := printSchema(page.Out, data, output); err != nil {
					return err
				}
			case "FORMAT_PROTOBUF":
				printProtoSchema(page.Out, data, filterPathPrefix, output)
			default:
				page.Stop()
				fmt.Printf("%s Unknown schema format: %s\n", term.Red(term.FailureIcon()), format)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&namespaceID, "namespace", "n", "", "provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")

	cmd.Flags().Int32VarP(&version, "version", "v", 0, "provide version number")

	cmd.Flags().StringVarP(&output, "output", "o", "", "the directory path to write the descriptor files, default is to print on stdout")

	cmd.Flags().StringVar(&filterPathPrefix, "filter-path", "", "filter protocol buffer files by path prefix, e.g., --filter-path=google/protobuf")

	return cmd
}

func printSchema(writer io.Writer, data []byte, output string) error {
	if output != "" {
		if err := os.WriteFile(output, data, 0666); err != nil {
			return err
		}
	}

	err := quick.Highlight(writer, string(data), "JSON", "terminal16m", "solarized-light")
	if err != nil {
		writer.Write(data)
	}
	return nil
}

func printProtoSchema(writer io.Writer, data []byte, filterPathPrefix string, output string) error {
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
		if filterPathPrefix != "" && !strings.HasPrefix(fdName, filterPathPrefix) {
			continue
		}
		filteredFds = append(filteredFds, fd)
	}

	protoPrinter := &protoprint.Printer{}

	if output != "" {
		if err := protoPrinter.PrintProtosToFileSystem(filteredFds, output); err != nil {
			return err
		}
	}

	var schema string

	for _, fd := range filteredFds {
		protoAsString, err := protoPrinter.PrintProtoToString(fd)
		if err != nil {
			return err
		}
		schema = schema + fmt.Sprintf("\n//Schema file:: %s\n\n%s", fd.GetName(), protoAsString)
	}

	err = quick.Highlight(writer, schema, "Protocol Buffer", "terminal16m", "solarized-light")
	if err != nil {
		fmt.Fprint(writer, schema)
	}

	return nil
}
