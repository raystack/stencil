package store

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errKind int

const (
	_ errKind = iota
	unknown
	conflict
	noRows
)

var (
	//UnknownErr default sentinel error for storage layer
	UnknownErr = StorageErr{kind: unknown}
	//ConflictErr can used for contraint violations
	ConflictErr = StorageErr{kind: conflict}
	//NoRowsErr can be used to represent not found/no result
	NoRowsErr = StorageErr{kind: noRows}
)

// StorageErr implements error interface. Used for storage layer. Consumers can check for this error type to inspect storage errors.
type StorageErr struct {
	name string
	kind errKind
	err  error
}

func (e StorageErr) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return e.name
}

// GRPCStatus this is used by gateway interceptor to return appropriate http status code and message
func (e StorageErr) GRPCStatus() *status.Status {
	if e.kind == noRows {
		return status.New(codes.NotFound, fmt.Sprintf("%s %s", e.name, "not found"))
	}
	if e.kind == conflict {
		return status.New(codes.AlreadyExists, fmt.Sprintf("%s %s", e.name, "resource already exists"))
	}
	return status.New(codes.Unknown, e.Error())
}

// WithErr convenience function to override sentinel errors
func (e StorageErr) WithErr(err error, name string) StorageErr {
	e.err = err
	e.name = name
	return e
}

func (e StorageErr) Unwrap() error {
	return e.err
}

func (e StorageErr) Is(err error) bool {
	sErr, ok := err.(StorageErr)
	if !ok {
		return false
	}
	return e.kind == sErr.kind
}
