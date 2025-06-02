-- +goose Up
CREATE TYPE role_type AS ENUM ('judge', 'participant', 'admin', 'organizer');

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name VARCHAR(128) NOT NULL,
    last_name VARCHAR(128) NOT NULL,
    surname VARCHAR(128) NOT NULL,
    phone_number VARCHAR(20),
    birth_date DATE NOT NULL,
    gender VARCHAR(10),
    role role_type NOT NULL
);

CREATE TABLE participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    ovz INTEGER NOT NULL,
    school_name VARCHAR(255) NOT NULL,
    city VARCHAR(128) NOT NULL,
    reason TEXT,
    citizenship VARCHAR(64),
    class_number INTEGER NOT NULL
);

CREATE INDEX idx_participants_user_id ON participants(user_id);


-- +goose Down
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS role_type;
