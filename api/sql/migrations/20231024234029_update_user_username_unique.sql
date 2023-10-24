-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX uq_user_username ON user(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uq_user_username ON user;
-- +goose StatementEnd
