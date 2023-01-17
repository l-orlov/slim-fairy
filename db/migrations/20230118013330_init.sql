-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL    DEFAULT '',
    age  int  NOT NULL    DEFAULT 0
);

COMMENT ON TABLE users IS 'Пользователи';
COMMENT ON COLUMN users.id IS 'ID пользователя';
COMMENT ON COLUMN users.name IS 'Имя пользователя';
COMMENT ON COLUMN users.age IS 'Возраст пользователя';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
