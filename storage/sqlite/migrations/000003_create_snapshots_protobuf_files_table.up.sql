CREATE TABLE IF NOT EXISTS snapshots_protobuf_files(
	snapshot_id INTEGER,
	file_id INTEGER,
	CONSTRAINT fk_snapshot FOREIGN KEY(snapshot_id) REFERENCES snapshots(id),
	CONSTRAINT fk_file FOREIGN KEY(file_id) REFERENCES protobuf_files(id)
);

CREATE UNIQUE INDEX snapshot_id_file_id_idx ON snapshots_protobuf_files (snapshot_id, file_id)
