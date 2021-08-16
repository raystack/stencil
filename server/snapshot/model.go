package snapshot

import "errors"

var (
	//ErrNotFound used when snapshot is not found
	ErrNotFound = errors.New("not found")
)

// Snapshot represents specific version of protodescriptorset
type Snapshot struct {
	ID        int64
	Namespace string
	Name      string
	Version   string
	Latest    bool
}
