package store

import (
	"fmt"

	"connectrpc.com/connect"
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

// ConnectError returns an appropriate Connect error for this storage error.
func (e StorageErr) ConnectError() *connect.Error {
	if e.kind == noRows {
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("%s %s", e.name, "not found"))
	}
	if e.kind == conflict {
		return connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("%s %s", e.name, "resource already exists"))
	}
	return connect.NewError(connect.CodeUnknown, e)
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

// Is reports whether the target error matches this StorageErr's kind.
func (e StorageErr) Is(err error) bool {
	sErr, ok := err.(StorageErr)
	if !ok {
		return false
	}
	return e.kind == sErr.kind
}
