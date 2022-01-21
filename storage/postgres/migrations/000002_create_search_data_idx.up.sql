CREATE INDEX search_data_idx ON schema_files USING gin (search_data);
CREATE INDEX search_data_fields_idx ON schema_files USING gin ((search_data->'Fields'));