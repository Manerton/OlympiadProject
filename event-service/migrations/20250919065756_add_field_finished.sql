-- +goose Up
ALTER TABLE events
    ADD COLUMN finished INT NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE events
    DROP COLUMN finished;
