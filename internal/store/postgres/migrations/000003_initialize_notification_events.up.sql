CREATE TABLE IF NOT EXISTS notification_events(
    id                  VARCHAR PRIMARY KEY,
    type                VARCHAR,
    event_time           TIMESTAMP,
    namespace_id        VARCHAR,
    schema_id           BIGINT,
    version_id          VARCHAR,
    success             BOOLEAN,
    created_at          TIMESTAMP,
    updated_at          TIMESTAMP
);