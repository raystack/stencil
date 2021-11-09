package validator

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func checkIfFieldRequired(f protoreflect.FieldDescriptor) bool {
	opts := f.Options()
	v := opts.ProtoReflect().Get(annotations.E_FieldBehavior.TypeDescriptor())
	if v.List().IsValid() {
		l := v.List()
		for i := 0; i < l.Len(); i++ {
			eVal := l.Get(i)
			if annotations.FieldBehavior(eVal.Enum()) == annotations.FieldBehavior_REQUIRED {
				return true
			}
		}
	}
	return false
}

func checkValueExists(kind protoreflect.Kind, v protoreflect.Value) bool {
	switch kind {
	case protoreflect.BytesKind:
		d := v.Bytes()
		return len(d) > 0
	case protoreflect.StringKind:
		d := v.String()
		return len(d) > 0
	case protoreflect.EnumKind:
		d := v.Enum()
		// This relies on the convention that, in enum definition 0 value is UNSPECIFIED
		return d != 0
	default:
		return true
	}
}

func addPrefix(prefix string, fieldNames []string) []string {
	for i := 0; i < len(fieldNames); i++ {
		field := fieldNames[i]
		fieldNames[i] = fmt.Sprintf("%s.%s", prefix, field)
	}
	return fieldNames
}

func validateMessage(m protoreflect.ProtoMessage) []string {
	var missingFields []string
	md := m.ProtoReflect().Descriptor()
	fds := md.Fields()
	for i := 0; i < fds.Len(); i++ {
		fd := fds.Get(i)
		v := m.ProtoReflect().Get(fd)
		if fd.Kind() == protoreflect.MessageKind && proto.Size(v.Message().Interface()) != 0 {
			nestedFields := validateMessage(v.Message().Interface())
			prefixedFields := addPrefix(fd.JSONName(), nestedFields)
			missingFields = append(missingFields, prefixedFields...)
		}
		if checkIfFieldRequired(fd) && !checkValueExists(fd.Kind(), v) {
			missingFields = append(missingFields, fd.JSONName())
		}
	}
	return missingFields
}

func validate(req interface{}) error {
	msg, ok := req.(proto.Message)
	if !ok {
		return nil
	}
	missingFields := validateMessage(msg)
	if len(missingFields) == 0 {
		return nil
	}
	return status.Error(codes.InvalidArgument, fmt.Sprintf("following fields are missing: %s", strings.Join(missingFields, ", ")))
}

// UnaryServerInterceptor interceptor for validating required fields
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if validErr := validate(req); validErr != nil {
			return nil, validErr
		}
		return handler(ctx, req)
	}
}
