CREATE TABLE IF NOT EXISTS namespaces(
	id VARCHAR PRIMARY KEY,
	format VARCHAR,
	compatibility VARCHAR,
	description VARCHAR,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS schemas(
	id VARCHAR PRIMARY KEY,
	authority VARCHAR,
	format VARCHAR,
	compatibility VARCHAR,
	description VARCHAR,
	namespace_id VARCHAR,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	CONSTRAINT fk_schemas_namespace_id FOREIGN KEY(namespace_id) REFERENCES namespaces(id)
);

CREATE TABLE IF NOT EXISTS versions(
	version BIGINT,
	schema_id VARCHAR,
	id VARCHAR,
	data bytea,
	created_at TIMESTAMP,
	CONSTRAINT fk_versions_schema_id FOREIGN KEY(schema_id) REFERENCES schemas(id),
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
	CONSTRAINT fk_versions_schema_files_version_id FOREIGN KEY(version_id) REFERENCES versions(id),
	CONSTRAINT fk_versions_schema_files_schema_file_id FOREIGN KEY(schema_file_id) REFERENCES schema_files(id)
);
