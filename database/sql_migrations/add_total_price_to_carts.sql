-- +migrate Up
ALTER TABLE carts ADD COLUMN IF NOT EXISTS total_price INT;

-- +migrate Down
ALTER TABLE carts DROP COLUMN IF EXISTS total_price;