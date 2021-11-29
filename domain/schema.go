package domain

import "context"

// Metadata model
type Metadata struct {
	Authority     string
	Format        string
	Compatibility string
}

type SchemaInfo struct {
	ID       string
	Version  int32
	Location string
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
	GetLatestSchema(context.Context, string, string) ([]byte, error)
	GetSchemaMetadata(context.Context, string, string) (*Metadata, error)
	UpdateSchemaMetadata(context.Context, string, string, *Metadata) (*Metadata, error)
	DeleteSchema(context.Context, string, string) error
	DeleteVersion(context.Context, string, string, int32) error
}

//SchemaService Service interface for schema management
type SchemaService interface {
	Create(context.Context, string, string, *Metadata, []byte) (SchemaInfo, error)
	List(context.Context, string) ([]string, error)
	Get(context.Context, string, string, int32) ([]byte, error)
	Delete(context.Context, string, string) error
	GetLatest(context.Context, string, string) ([]byte, error)
	ListVersions(context.Context, string, string) ([]int32, error)
	GetMetadata(context.Context, string, string) (*Metadata, error)
	UpdateMetadata(context.Context, string, string, *Metadata) (*Metadata, error)
	DeleteVersion(context.Context, string, string, int32) error
}
