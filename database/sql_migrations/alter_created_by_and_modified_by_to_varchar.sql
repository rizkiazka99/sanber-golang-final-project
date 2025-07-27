-- +migrate Up
ALTER TABLE items
ALTER COLUMN created_by TYPE VARCHAR USING created_by::VARCHAR;

ALTER TABLE items
ALTER COLUMN modified_by TYPE VARCHAR USING modified_by::VARCHAR;

-- +migrate Down
ALTER TABLE items
ALTER COLUMN created_by TYPE TIMESTAMP USING created_by::TIMESTAMP;

ALTER TABLE items
ALTER COLUMN modified_by TYPE TIMESTAMP USING modified_by::TIMESTAMP;