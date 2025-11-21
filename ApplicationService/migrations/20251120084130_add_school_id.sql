-- +goose Up
-- +goose StatementBegin
ALTER TABLE applications
    ADD COLUMN school_id UUID DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE applications
DROP COLUMN school_id;
-- +goose StatementEnd
