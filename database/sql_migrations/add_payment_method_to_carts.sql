-- +migrate Up
ALTER TABLE carts ADD COLUMN IF NOT EXISTS payment_method VARCHAR(255);

-- +migrate Down
ALTER TABLE carts DROP COLUMN IF EXISTS payment_method;