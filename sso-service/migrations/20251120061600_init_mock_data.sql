-- +goose Up
INSERT INTO users (id, email, password_hash, firstname, surname, patronymic, phone_number, birthdate, gender, role, activated, created_at) VALUES
('33333333-3333-3333-3333-333333333331', 'admin@mail.ru',     '$2a$10$qu78o4axVYWAznSw3a3q8OErhUBbTmpgvWGtLgcLX.FhYAS9Q6mkW', 'Админ', 'Системный', 'Админович',  '+70000000001', '1980-01-01', 1, 1, true, now()),
('33333333-3333-3333-3333-333333333332', 'jury@mail.ru',      '$2a$10$qu78o4axVYWAznSw3a3q8OErhUBbTmpgvWGtLgcLX.FhYAS9Q6mkW', 'Иван',  'Судья',     'Петрович',   '+70000000002', '1985-02-02', 1, 3, true, now()),
('33333333-3333-3333-3333-333333333333', 'organizer@mail.ru', '$2a$10$qu78o4axVYWAznSw3a3q8OErhUBbTmpgvWGtLgcLX.FhYAS9Q6mkW', 'Ольга', 'Организатор','Сергеевна', '+70000000003', '1990-03-03', 2, 4, true, now());



-- +goose Down
DELETE FROM users
WHERE id IN (
    '33333333-3333-3333-3333-333333333331',
    '33333333-3333-3333-3333-333333333332',
    '33333333-3333-3333-3333-333333333333'
);