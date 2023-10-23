# name: GetUser :one
SELECT * FROM user
WHERE id = ? LIMIT 1;

# name: SignUpUser :execlastid
INSERT INTO user (username, password)
VALUES (?, ?);

# name: SetUserRefreshToken :exec
UPDATE user SET refresh_token = ?
WHERE id = ?;

# name: GetUserByUsername :one
SELECT * FROM user
WHERE username = ? LIMIT 1;
