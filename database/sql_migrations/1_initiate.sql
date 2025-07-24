-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE users (
    id BIGINT PRIMARY KEY NOT NULL,
    username varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    token TEXT,
    expire_time TIMESTAMP
);

CREATE TABLE items (
    id BIGINT PRIMARY KEY NOT NULL,
    item_name varchar(500) NOT NULL,
    description TEXT,
    price INT NOT NULL,
    stock INT NOT NULL,
    created_by TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    modified_by TIMESTAMP NOT NULL,
    modified_at TIMESTAMP NOT NULL
);

CREATE TABLE items_images (
    id BIGINT PRIMARY KEY NOT NULL,
    item_id BIGINT NOT NULL REFERENCES items(id),
    image_url TEXT NOT NULL
);

CREATE TABLE carts (
    id BIGINT PRIMARY KEY NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE cart_items (
    id BIGINT PRIMARY KEY NOT NULL,
    cart_id BIGINT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    item_id BIGINT NOT NULL REFERENCES items(id),
    quantity INT NOT NULL CHECK (quantity > 0)
);

CREATE TABLE transactions (
    id BIGINT PRIMARY KEY NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),
    cart_id BIGINT NOT NULL REFERENCES carts(id),
    total_price INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +migrate StatementEnd