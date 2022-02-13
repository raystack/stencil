package domain

import "context"

// Metadata model
type Metadata struct {
	Authority     string
	Format        string
	Compatibility string
}

type SchemaInfo struct {
	ID       string `json:"id"`
	Version  int32  `json:"version"`
	Location string `json:"location"`
}

type SchemaFile struct {
	ID     string
	Types  []string
	Fields []string
	Data   []byte
}

// SchemaRepository for Schema
type SchemaRepository interface {
	CreateSchema(ctx context.Context, namespace string, schema string, metadata *Metadata, versionID string, schemaFile *SchemaFile) (version int32, err error)
	ListSchemas(context.Context, string) ([]string, error)
	ListVersions(context.Context, string, string) ([]int32, error)
	GetSchema(context.Context, string, string, int32) ([]byte, error)
	GetLatestVersion(context.Context, string, string) (int32, error)
	GetSchemaMetadata(context.Context, string, string) (*Metadata, error)
	UpdateSchemaMetadata(context.Context, string, string, *Metadata) (*Metadata, error)
	DeleteSchema(context.Context, string, string) error
	DeleteVersion(context.Context, string, string, int32) error
}

//SchemaService Service interface for schema management
type SchemaService interface {
	Create(context.Context, string, string, *Metadata, []byte) (SchemaInfo, error)
	List(context.Context, string) ([]string, error)
	Get(context.Context, string, string, int32) (*Metadata, []byte, error)
	Delete(context.Context, string, string) error
	GetLatest(context.Context, string, string) (*Metadata, []byte, error)
	ListVersions(context.Context, string, string) ([]int32, error)
	GetMetadata(context.Context, string, string) (*Metadata, error)
	UpdateMetadata(context.Context, string, string, *Metadata) (*Metadata, error)
	DeleteVersion(context.Context, string, string, int32) error
}
