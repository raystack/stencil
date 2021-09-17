package visualize

import (
	"fmt"

	"github.com/emicklei/dot"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	NodeShape = "note"
	NodeStyle = "filled"
	NodeColor = "cornsilk"
)

func GetProtoFileDependencyGraph(file *descriptorpb.FileDescriptorSet) (*dot.Graph, error) {
	files, err := protodesc.NewFiles(file)
	if err != nil {
		return nil, err
	}

	di := dot.NewGraph(dot.Directed)
	di.Attr("rankdir", "LR")

	files.RangeFiles(func(file protoreflect.FileDescriptor) bool {
		sourceNode := di.Node(fmt.Sprintf("%s\n%s", string(file.Package()), file.Path()))
		setDefaultAttributes(sourceNode)

		buildGraph(di, sourceNode, file)
		return true
	})

	return di, nil
}

func buildGraph(di *dot.Graph, sourceNode dot.Node, file protoreflect.FileDescriptor) {
	for i := 0; i < file.Imports().Len(); i++ {
		imp := file.Imports().Get(i)
		destNode := di.Node(fmt.Sprintf("%s\n%s", string(imp.Package()), imp.Path()))
		setDefaultAttributes(destNode)

		di.Edge(sourceNode, destNode, "")
	}
}

func setDefaultAttributes(n dot.Node) {
	n.Attr("shape", NodeShape)
	n.Attr("style", NodeStyle)
	n.Attr("fillcolor", NodeColor)
}
