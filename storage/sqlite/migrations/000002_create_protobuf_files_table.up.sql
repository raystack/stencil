CREATE TABLE IF NOT EXISTS protobuf_files(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	search_data JSONB,
	data bytea
);
CREATE UNIQUE INDEX data_unique_idx ON protobuf_files (data);
