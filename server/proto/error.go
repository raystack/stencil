package proto

import (
	"fmt"
	"strings"

	"go.uber.org/multierr"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type validationErr struct {
	desc protoreflect.Descriptor
	msgs []string
}

func (v *validationErr) add(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	v.msgs = append(v.msgs, msg)
}

func (v *validationErr) isNil() bool {
	return len(v.msgs) == 0
}

func (v *validationErr) Error() string {
	var msg strings.Builder
	file := getFileDescriptor(v.desc)
	path := file.Path()
	for i := 0; i < len(v.msgs); i++ {
		fmt.Fprintf(&msg, "%s: %s; ", path, v.msgs[i])
	}
	return strings.TrimSuffix(msg.String(), "; ")
}

func newValidationErr(file protoreflect.Descriptor) *validationErr {
	return &validationErr{
		desc: file,
	}
}

func combineErr(err error, validationErr *validationErr) error {
	if validationErr.isNil() {
		return err
	}
	return multierr.Combine(err, validationErr)
}
