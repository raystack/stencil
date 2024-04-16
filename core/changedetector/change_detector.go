package changedetector

import "time"

type ChangeRequest struct {
	NamespaceID string
	SchemaName  string
	Version     int32
	VersionID   string
	OldData     []byte
	NewData     []byte
}

type NotificationEvent struct {
	ID          string
	Type        string
	EventTime   time.Time
	NamespaceID string
	SchemaID    int32
	VersionID   string
	Success     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
