-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    event_id UUID NOT NULL,
    status INT NOT NULL DEFAULT 1, -- 1 = не обработано, 2 = одобрено, 3 = отклонено
    reason INT DEFAULT NULL,       -- 1 = по результатам предыдущего года, 2 = по результатам текущего
    code TEXT DEFAULT NULL,        -- например: 09_111_25
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()

    CONSTRAINT unique_user_event UNIQUE (user_id, event_id)
);

-- +goose Down
DROP TABLE IF EXISTS applications;
