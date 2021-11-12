CREATE TABLE IF NOT EXISTS namespaces(
	id VARCHAR PRIMARY KEY,
	format VARCHAR,
	compatibility VARCHAR,
	description VARCHAR,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS schemas(
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR NOT NULL,
	authority VARCHAR,
	format VARCHAR,
	compatibility VARCHAR,
	description VARCHAR,
	namespace_id VARCHAR NOT NULL,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	CONSTRAINT fk_schemas_namespace_id FOREIGN KEY(namespace_id) REFERENCES namespaces(id) ON DELETE CASCADE,
	CONSTRAINT schema_name_namespace_unique_idx UNIQUE (name, namespace_id)
);

CREATE TABLE IF NOT EXISTS versions(
	id VARCHAR,
	version BIGINT,
	schema_id BIGINT,
	created_at TIMESTAMP,
	CONSTRAINT fk_versions_schema_id FOREIGN KEY(schema_id) REFERENCES schemas(id) ON DELETE CASCADE,
	CONSTRAINT schema_version_unique_idx UNIQUE (version, schema_id),
	CONSTRAINT schema_id_unique UNIQUE (id)
);

CREATE TABLE IF NOT EXISTS schema_files(
	id VARCHAR,
	search_data JSONB,
	data bytea,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	CONSTRAINT schema_files_id_unique_idx UNIQUE (id)
);

CREATE TABLE IF NOT EXISTS versions_schema_files(
	version_id VARCHAR,
	schema_file_id VARCHAR,
	CONSTRAINT fk_versions_schema_files_version_id FOREIGN KEY(version_id) REFERENCES versions(id) ON DELETE CASCADE,
	CONSTRAINT fk_versions_schema_files_schema_file_id FOREIGN KEY(schema_file_id) REFERENCES schema_files(id)
);
