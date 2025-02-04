CREATE TABLE container_status (
    id SERIAL PRIMARY KEY,
    ip_address INET NOT NULL UNIQUE,
    ping_time NUMERIC,
    last_successful_ping TIMESTAMP,
    updated_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now()
);
