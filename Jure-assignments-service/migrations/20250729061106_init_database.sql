-- +goose Up
-- +goose StatementBegin
CREATE TABLE jury_assignments ( 
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    event_id UUID NOT NULL
);

CREATE INDEX idx_event_jury ON jury_assignments(event_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS jury_assignments;
-- +goose StatementEnd
