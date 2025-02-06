CREATE INDEX idx_last_successful_ping ON container_status(last_successful_ping);

CREATE INDEX idx_updated_at ON container_status(updated_at);