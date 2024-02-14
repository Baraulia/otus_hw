-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    header VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    user_id VARCHAR(255) NOT NULL,
    event_time TIMESTAMPTZ NOT NULL,
    finish_event_time TIMESTAMPTZ,
    notification_time TIMESTAMPTZ
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event;
-- +goose StatementEnd
