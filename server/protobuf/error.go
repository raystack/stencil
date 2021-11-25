package protobuf

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type diff struct {
	kind diffKind
	msg  string
}

type compatibilityErr struct {
	notAllowed []diffKind
	diffs      []diff
}

func (c *compatibilityErr) add(kind diffKind, desc protoreflect.Descriptor, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	path := desc.ParentFile().Path()
	if kind.contains(c.notAllowed) && msg != "" {
		c.diffs = append(c.diffs, diff{kind: kind, msg: fmt.Sprintf("%s: %s", path, msg)})
	}
}

func (c *compatibilityErr) addDiffs(diffs []diff) {
	c.diffs = append(c.diffs, diffs...)
}

func (c *compatibilityErr) isEmpty() bool {
	return len(c.diffs) > 0
}

func (c *compatibilityErr) Error() string {
	var msgs []string
	for _, val := range c.diffs {
		msgs = append(msgs, val.msg)
	}
	return strings.Join(msgs, ";")
}

func (c *compatibilityErr) GRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, c.Error())
}
