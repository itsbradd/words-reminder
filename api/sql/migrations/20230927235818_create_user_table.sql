-- +goose Up
-- +goose StatementBegin
-- CREATE TABLE IF NOT EXISTS user (
--     id INT
-- )
CREATE TABLE IF NOT EXISTS user (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(1000)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
-- +goose StatementEnd
