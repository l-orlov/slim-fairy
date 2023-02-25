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

CREATE TABLE clients
(
    id         uuid PRIMARY KEY     DEFAULT gen_random_uuid(),
    name       text        NOT NULL DEFAULT '',
    email      text        NOT NULL DEFAULT '',
    phone      text        NOT NULL DEFAULT '',
    age        int         NOT NULL DEFAULT 0,
    weight     int         NOT NULL DEFAULT 0,
    height     int         NOT NULL DEFAULT 0,
    gender     int         NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

COMMENT ON TABLE clients IS 'Clients';
COMMENT ON COLUMN clients.id IS 'Client ID';
COMMENT ON COLUMN clients.name IS 'Client Name';
COMMENT ON COLUMN clients.email IS 'Client Email';
COMMENT ON COLUMN clients.phone IS 'Client Phone';
COMMENT ON COLUMN clients.age IS 'Client Age';
COMMENT ON COLUMN clients.weight IS 'Client Weight';
COMMENT ON COLUMN clients.height IS 'Client Height';
COMMENT ON COLUMN clients.gender IS 'Client Gender';
COMMENT ON COLUMN clients.created_at IS 'Create date';
COMMENT ON COLUMN clients.updated_at IS 'Update date';

CREATE TRIGGER update_clients_updated_at
    BEFORE UPDATE
    ON clients
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column();

CREATE UNIQUE INDEX ON clients USING btree (email);

CREATE TABLE nutritionists
(
    id         uuid PRIMARY KEY     DEFAULT gen_random_uuid(),
    name       text        NOT NULL DEFAULT '',
    email      text        NOT NULL DEFAULT '',
    phone      text        NOT NULL DEFAULT '',
    age        int         NOT NULL DEFAULT 0,
    gender     int         NOT NULL DEFAULT 0,
    info       text        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT (now() at time zone 'utc')
);

COMMENT ON TABLE nutritionists IS 'Nutritionists';
COMMENT ON COLUMN nutritionists.id IS 'Nutritionist ID';
COMMENT ON COLUMN nutritionists.name IS 'Nutritionist Name';
COMMENT ON COLUMN nutritionists.email IS 'Nutritionist Email';
COMMENT ON COLUMN nutritionists.phone IS 'Nutritionist Phone';
COMMENT ON COLUMN nutritionists.age IS 'Nutritionist Age';
COMMENT ON COLUMN nutritionists.gender IS 'Nutritionist Gender';
COMMENT ON COLUMN nutritionists.info IS 'Nutritionist Info';
COMMENT ON COLUMN nutritionists.created_at IS 'Create date';
COMMENT ON COLUMN nutritionists.updated_at IS 'Update date';

CREATE TRIGGER update_nutritionists_updated_at
    BEFORE UPDATE
    ON nutritionists
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column();

CREATE UNIQUE INDEX ON nutritionists USING btree (email);

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
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS nutritionists;
DROP TABLE IF EXISTS auth_data;
-- +goose StatementEnd
