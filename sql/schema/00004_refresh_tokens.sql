-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE refresh_tokens
(
    token      TEXT PRIMARY KEY,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id    UUID         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMP    NOT NULL,
    revoked_at TIMESTAMP


);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE refresh_tokens;
-- +goose StatementEnd