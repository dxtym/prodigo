CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

INSERT INTO users (username, role, password) 
VALUES ('admin', 'admin', '$2a$10$ZOE2qCYDGSk9yeDfykQGzerHCAqRFs8/mVEukM6jU9FHr8MnAt/bK')
ON CONFLICT (username) DO NOTHING;