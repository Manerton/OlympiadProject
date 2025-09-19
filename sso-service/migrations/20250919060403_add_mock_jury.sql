-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, email, password_hash, firstname, surname, patronymic, phone_number, birthdate, gender, role, activated, created_at) VALUES
('33333333-3333-3333-3333-33333333333b', 'jury2@mail.ru', '$2a$10$qu78o4axVYWAznSw3a3q8OErhUBbTmpgvWGtLgcLX.FhYAS9Q6mkW',
 'Павел', 'Громов', 'Викторович', '+70000000011', '1982-11-11', 1, 3, true, now()),
('33333333-3333-3333-3333-33333333333c', 'jury3@mail.ru', '$2a$10$qu78o4axVYWAznSw3a3q8OErhUBbTmpgvWGtLgcLX.FhYAS9Q6mkW',
 'Светлана', 'Кириллова', 'Игоревна', '+70000000012', '1987-12-12', 2, 3, true, now()),
('33333333-3333-3333-3333-33333333333d', 'jury4@mail.ru', '$2a$10$qu78o4axVYWAznSw3a3q8OErhUBbTmpgvWGtLgcLX.FhYAS9Q6mkW',
 'Михаил', 'Савельев', 'Алексеевич', '+70000000013', '1979-09-15', 1, 3, true, now()),
('33333333-3333-3333-3333-33333333333e', 'jury5@mail.ru', '$2a$10$qu78o4axVYWAznSw3a3q8OErhUBbTmpgvWGtLgcLX.FhYAS9Q6mkW',
 'Татьяна', 'Власова', 'Сергеевна', '+70000000014', '1992-04-20', 2, 3, true, now());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users
WHERE id IN (
    '33333333-3333-3333-3333-33333333333b',
    '33333333-3333-3333-3333-33333333333c',
    '33333333-3333-3333-3333-33333333333d',
    '33333333-3333-3333-3333-33333333333e'
);
-- +goose StatementEnd
