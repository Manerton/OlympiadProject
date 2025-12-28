-- +goose Up
ALTER TABLE events
    RENAME COLUMN finished TO status;

-- +goose Down
ALTER TABLE events
    RENAME COLUMN status TO finished;
