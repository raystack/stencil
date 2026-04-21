package schema

import "context"

// Metadata represents schema metadata such as format and compatibility.
type Metadata struct {
	Authority     string
	Format        string
	Compatibility string
}

// SchemaInfo represents version and location information for a schema.
type SchemaInfo struct {
	ID       string `json:"id"`
	Version  int32  `json:"version"`
	Location string `json:"location"`
}

// SchemaFile represents a schema's data along with its extracted types and fields.
type SchemaFile struct {
	ID     string
	Types  []string
	Fields []string
	Data   []byte
}

// Repository defines the persistence interface for schemas.
type Repository interface {
	Create(ctx context.Context, namespace string, schema string, metadata *Metadata, versionID string, schemaFile *SchemaFile) (version int32, err error)
	List(context.Context, string) ([]Schema, error)
	ListVersions(context.Context, string, string) ([]int32, error)
	Get(context.Context, string, string, int32) ([]byte, error)
	GetLatestVersion(context.Context, string, string) (int32, error)
	GetMetadata(context.Context, string, string) (*Metadata, error)
	UpdateMetadata(context.Context, string, string, *Metadata) (*Metadata, error)
	Delete(context.Context, string, string) error
	DeleteVersion(context.Context, string, string, int32) error
}

// ParsedSchema defines the interface for a parsed schema with compatibility checks.
type ParsedSchema interface {
	IsBackwardCompatible(ParsedSchema) error
	IsForwardCompatible(ParsedSchema) error
	IsFullCompatible(ParsedSchema) error
	Format() string
	GetCanonicalValue() *SchemaFile
}

// Provider defines the interface for parsing raw schema data into a ParsedSchema.
type Provider interface {
	ParseSchema(format string, data []byte) (ParsedSchema, error)
}

// Cache defines a key-value cache interface for storing schema data.
type Cache interface {
	Get(interface{}) (interface{}, bool)
	Set(interface{}, interface{}, int64) bool
	Del(interface{})
}

// Schema represents a schema entity with its configuration.
type Schema struct {
	Name          string
	Format        string
	Compatibility string
	Authority     string
}
