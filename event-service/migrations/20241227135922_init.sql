-- +goose Up
-- +goose StatementBegin
CREATE TYPE event_type AS ENUM ('REGIONAL_STAGE', 'OLYMPIAD', 'STAGE', 'VIEW_WORKS', 'APPEAL');

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    event_type event_type NOT NULL,
    previous_event_id INT,
    subject VARCHAR(128),
    additional_info TEXT,
    CONSTRAINT fk_previous_event
        FOREIGN KEY (previous_event_id)
        REFERENCES events (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
