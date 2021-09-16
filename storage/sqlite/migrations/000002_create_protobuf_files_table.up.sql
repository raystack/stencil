CREATE TABLE IF NOT EXISTS protobuf_files(
	id BIGSERIAL PRIMARY KEY,
	search_data JSONB,
	data bytea
);
CREATE UNIQUE INDEX data_unique_idx ON protobuf_files (data);
