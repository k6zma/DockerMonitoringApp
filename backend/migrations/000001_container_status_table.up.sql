CREATE TABLE container_status (
    id BIGSERIAL PRIMARY KEY,
    ip_address INET NOT NULL UNIQUE,
    ping_time DOUBLE PRECISION NOT NULL,
    last_successful_ping TIMESTAMP,
    updated_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now()
);
