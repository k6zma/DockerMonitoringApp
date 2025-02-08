CREATE TABLE container_status (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT '',
    ip_address INET NOT NULL UNIQUE,
    status VARCHAR(255) NOT NULL DEFAULT 'created',
    ping_time DOUBLE PRECISION NULL,
    last_successful_ping TIMESTAMP,
    updated_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now()
);
