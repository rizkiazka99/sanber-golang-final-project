-- +migrate Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(255) NOT NULL;

-- +migrate Down
ALTER TABLE users DROP COLUMN IF EXISTS role;