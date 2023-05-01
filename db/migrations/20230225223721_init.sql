-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION set_updated_at_column() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now() at time zone 'utc';
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE users
(
    id                      uuid PRIMARY KEY     DEFAULT gen_random_uuid(),
    name                    text        NOT NULL DEFAULT '',
    email                   text,
    phone                   text,
    telegram_id             text,
    age                     int,
    weight                  int,
    height                  int,
    gender                  int,
    physical_activity_level int,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

COMMENT ON TABLE users IS 'Users';
COMMENT ON COLUMN users.id IS 'ID';
COMMENT ON COLUMN users.name IS 'Name';
COMMENT ON COLUMN users.email IS 'Email';
COMMENT ON COLUMN users.phone IS 'Phone';
COMMENT ON COLUMN users.telegram_id IS 'Telegram ID';
COMMENT ON COLUMN users.age IS 'Age';
COMMENT ON COLUMN users.weight IS 'Weight';
COMMENT ON COLUMN users.height IS 'Height';
COMMENT ON COLUMN users.gender IS 'Gender';
COMMENT ON COLUMN users.physical_activity_level IS 'Physical activity level';
COMMENT ON COLUMN users.created_at IS 'Create date';
COMMENT ON COLUMN users.updated_at IS 'Update date';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column();

CREATE UNIQUE INDEX ON users USING btree (email);
CREATE UNIQUE INDEX ON users USING btree (phone);
CREATE UNIQUE INDEX ON users USING btree (telegram_id);

CREATE TABLE nutritionists
(
    id          uuid PRIMARY KEY     DEFAULT gen_random_uuid(),
    name        text        NOT NULL DEFAULT '',
    email       text,
    phone       text,
    telegram_id text,
    age         int,
    gender      int,
    info        text,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

COMMENT ON TABLE nutritionists IS 'Nutritionists';
COMMENT ON COLUMN nutritionists.id IS 'ID';
COMMENT ON COLUMN nutritionists.name IS 'Name';
COMMENT ON COLUMN nutritionists.email IS 'Email';
COMMENT ON COLUMN nutritionists.phone IS 'Phone';
COMMENT ON COLUMN users.telegram_id IS 'Telegram ID';
COMMENT ON COLUMN nutritionists.age IS 'Age';
COMMENT ON COLUMN nutritionists.gender IS 'Gender';
COMMENT ON COLUMN nutritionists.info IS 'Info';
COMMENT ON COLUMN nutritionists.created_at IS 'Create date';
COMMENT ON COLUMN nutritionists.updated_at IS 'Update date';

CREATE TRIGGER update_nutritionists_updated_at
    BEFORE UPDATE
    ON nutritionists
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column();

CREATE UNIQUE INDEX ON nutritionists USING btree (email);
CREATE UNIQUE INDEX ON nutritionists USING btree (phone);
CREATE UNIQUE INDEX ON nutritionists USING btree (telegram_id);

CREATE TABLE auth_data
(
    source_id   uuid        NOT NULL,
    source_type text        NOT NULL,
    password    text        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

COMMENT ON TABLE auth_data IS 'Authentication data';
COMMENT ON COLUMN auth_data.source_id IS 'Source ID';
COMMENT ON COLUMN auth_data.source_type IS 'Source type';
COMMENT ON COLUMN auth_data.password IS 'Password';
COMMENT ON COLUMN auth_data.created_at IS 'Create date';
COMMENT ON COLUMN auth_data.updated_at IS 'Update date';

CREATE TRIGGER update_auth_data_updated_at
    BEFORE UPDATE
    ON auth_data
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column();

CREATE UNIQUE INDEX ON auth_data USING btree (source_id, source_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS nutritionists;
DROP TABLE IF EXISTS auth_data;
-- +goose StatementEnd
