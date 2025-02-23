-- +goose Up
-- +goose StatementBegin
CREATE TYPE event_type AS ENUM ('REGIONAL_STAGE', 'OLYMPIAD', 'STAGE', 'VIEW_WORKS', 'APPEAL');

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    event_type event_type NOT NULL,
    previous_event_id UUID,
    subject VARCHAR(128),
    additional_info TEXT,
    CONSTRAINT fk_previous_event 
        FOREIGN KEY (previous_event_id) 
        REFERENCES events(id) 
        ON UPDATE CASCADE 
        ON DELETE RESTRICT
);

-- Добавляем индекс для поля previous_event_id
CREATE INDEX idx_previous_event ON events(previous_event_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
