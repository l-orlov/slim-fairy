-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION set_updated_at_column() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now() AT TIME ZONE 'utc';
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE users
(
    id                      UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    name                    TEXT        NOT NULL DEFAULT '',
    email                   TEXT,
    phone                   TEXT,
    telegram_id             BIGINT,
    age                     INT,
    weight                  INT,
    height                  INT,
    gender                  INT,
    physical_activity_level INT,
    created_by              TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
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
COMMENT ON COLUMN users.created_by IS 'Created by';
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
    id          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL DEFAULT '',
    email       TEXT,
    phone       TEXT,
    telegram_id BIGINT,
    age         INT,
    gender      INT,
    info        TEXT,
    created_by  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
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
COMMENT ON COLUMN users.created_by IS 'Created by';
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
    source_id   UUID        NOT NULL,
    source_type TEXT        NOT NULL,
    password    TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
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

CREATE TABLE chat_bot_dialogs
(
    id               UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_telegram_id BIGINT      NOT NULL,
    kind             TEXT        NOT NULL,
    status           INT         NOT NULL,
    data             JSONB,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

COMMENT ON TABLE chat_bot_dialogs IS 'Chat-bot dialogs with users';
COMMENT ON COLUMN chat_bot_dialogs.id IS 'ID';
COMMENT ON COLUMN chat_bot_dialogs.user_telegram_id IS 'User Telegram ID';
COMMENT ON COLUMN chat_bot_dialogs.kind IS 'Kind';
COMMENT ON COLUMN chat_bot_dialogs.status IS 'Status';
COMMENT ON COLUMN chat_bot_dialogs.data IS 'Dialog data';
COMMENT ON COLUMN nutritionists.created_at IS 'Create date';
COMMENT ON COLUMN nutritionists.updated_at IS 'Update date';

CREATE TRIGGER chat_bot_dialogs_updated_at
    BEFORE UPDATE
    ON chat_bot_dialogs
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column();

CREATE INDEX ON chat_bot_dialogs USING btree (user_telegram_id, kind, status);

CREATE TABLE ai_api_logs
(
    id          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    prompt      TEXT        NOT NULL,
    response    TEXT,
    user_id     UUID,
    source_id   UUID,
    source_type TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

COMMENT ON TABLE ai_api_logs IS 'AI API logs';
COMMENT ON COLUMN ai_api_logs.id IS 'ID';
COMMENT ON COLUMN ai_api_logs.prompt IS 'Prompt';
COMMENT ON COLUMN ai_api_logs.response IS 'Response';
COMMENT ON COLUMN ai_api_logs.user_id IS 'User ID';
COMMENT ON COLUMN ai_api_logs.source_id IS 'Source ID';
COMMENT ON COLUMN ai_api_logs.source_type IS 'Source Type';
COMMENT ON COLUMN ai_api_logs.created_at IS 'Create date';
COMMENT ON COLUMN ai_api_logs.updated_at IS 'Update date';

CREATE TRIGGER update_ai_api_logs_updated_at
    BEFORE UPDATE
    ON ai_api_logs
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column();

CREATE INDEX ON ai_api_logs USING btree (source_id, source_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS nutritionists;
DROP TABLE IF EXISTS auth_data;
-- +goose StatementEnd
