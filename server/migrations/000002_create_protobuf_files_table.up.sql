CREATE TABLE IF NOT EXISTS protobuf_files(
	id BIGSERIAL PRIMARY KEY,
	search_data JSONB,
	data bytea
);
CREATE INDEX search_data_idx ON protobuf_files USING gin (search_data);
CREATE INDEX search_data_messages_idx ON protobuf_files USING gin ((search_data->'messages'));
CREATE UNIQUE INDEX data_unique_idx ON protobuf_files (md5(data));
