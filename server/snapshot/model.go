package snapshot

import (
	"errors"
)

var (
	//ErrNotFound used when snapshot is not found
	ErrNotFound = errors.New("not found")
)

// Snapshot represents specific version of protodescriptorset
type Snapshot struct {
	ID        int64  `validate:"required"`
	Namespace string `validate:"required"`
	Name      string `validate:"required"`
	Version   string `validate:"required,version"`
	Latest    bool   `validate:"required"`
}
