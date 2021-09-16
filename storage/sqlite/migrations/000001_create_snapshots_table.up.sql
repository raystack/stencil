CREATE TABLE IF NOT EXISTS snapshots(
	id BIGSERIAL PRIMARY KEY,
	namespace VARCHAR NOT NULL,
	name VARCHAR NOT NULL,
	version VARCHAR NOT NULL,
	latest BOOLEAN,
	UNIQUE(namespace, name, version)
);
CREATE UNIQUE INDEX latest_unique_per_name_idx ON snapshots (namespace, name, latest)
WHERE latest = true;
