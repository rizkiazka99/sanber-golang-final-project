-- +migrate Up
DROP TABLE IF EXISTS transactions;

-- +migrate Down
CREATE TABLE transactions (
    id BIGINT PRIMARY KEY NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),
    cart_id BIGINT NOT NULL REFERENCES carts(id),
    total_price INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);